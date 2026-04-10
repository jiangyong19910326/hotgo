// Package plc 西门子 S7 PLC 通信驱动
package plc

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/robinson/gos7"
)

var helper = gos7.Helper{}

// Client S7 TCP 连接客户端，线程安全，支持自动重连
type Client struct {
	host     string
	port     int
	rack     int
	slot     int
	handler  *gos7.TCPClientHandler
	s7client gos7.Client
	mu       sync.Mutex
	lastPing time.Time
}

// NewClient 创建新的 S7 客户端（不立即连接）
func NewClient(host string, port, rack, slot int) *Client {
	return &Client{
		host: host,
		port: port,
		rack: rack,
		slot: slot,
	}
}

// connect 建立连接（调用方需持有锁）
func (c *Client) connect() error {
	h := gos7.NewTCPClientHandler(
		fmt.Sprintf("%s:%d", c.host, c.port),
		c.rack,
		c.slot,
	)
	h.Timeout = 5 * time.Second
	h.IdleTimeout = 10 * time.Second

	if err := h.Connect(); err != nil {
		return fmt.Errorf("plc connect %s:%d: %w", c.host, c.port, err)
	}

	c.handler = h
	c.s7client = gos7.NewClient(h)
	c.lastPing = time.Now()
	return nil
}

// ensureConnected 确保连接可用，断线则重连
func (c *Client) ensureConnected() error {
	if c.handler == nil {
		return c.connect()
	}
	// 超过 30s 未使用则主动重连，防止 PLC 侧超时断开
	if time.Since(c.lastPing) > 30*time.Second {
		_ = c.handler.Close()
		c.handler = nil
		c.s7client = nil
		return c.connect()
	}
	return nil
}

// Close 关闭连接
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.handler != nil {
		_ = c.handler.Close()
		c.handler = nil
		c.s7client = nil
	}
}

// ReadAll 批量读取数据点，返回每个点的结果。
// 读取失败的点 Valid=false，不影响其他点继续读取。
func (c *Client) ReadAll(reqs []ReadRequest) []ReadResult {
	c.mu.Lock()
	defer c.mu.Unlock()

	results := make([]ReadResult, len(reqs))

	if err := c.ensureConnected(); err != nil {
		for i, req := range reqs {
			results[i] = ReadResult{PointID: req.PointID, Field: req.Field, Valid: false}
		}
		return results
	}

	c.lastPing = time.Now()

	for i, req := range reqs {
		raw, rawStr, err := c.readSingle(req)
		if err != nil {
			results[i] = ReadResult{PointID: req.PointID, Field: req.Field, Valid: false}
			continue
		}

		engValue := raw*req.Scale + req.OffsetVal
		engValue = math.Round(engValue*10000) / 10000

		alarmType := 0
		if req.AlarmMax != nil && engValue > *req.AlarmMax {
			alarmType = 1
		} else if req.AlarmMin != nil && engValue < *req.AlarmMin {
			alarmType = 2
		}

		results[i] = ReadResult{
			PointID:   req.PointID,
			Field:     req.Field,
			RawValue:  rawStr,
			EngValue:  engValue,
			Valid:     true,
			AlarmType: alarmType,
			AlarmMin:  req.AlarmMin,
			AlarmMax:  req.AlarmMax,
			Unit:      req.Unit,
		}
	}
	return results
}

// readSingle 读取单个数据点，返回 (原始float64, 原始字符串, error)
func (c *Client) readSingle(req ReadRequest) (float64, string, error) {
	switch req.DataType {
	case TypeBool:
		buf := make([]byte, 1)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		bit := (buf[0] >> uint(req.BitOffset)) & 0x01
		v := float64(bit)
		return v, strconv.FormatFloat(v, 'f', 0, 64), nil

	case TypeByte:
		buf := make([]byte, 1)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		v := float64(buf[0])
		return v, strconv.FormatFloat(v, 'f', 0, 64), nil

	case TypeWord:
		buf := make([]byte, 2)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		var val uint16
		helper.GetValueAt(buf, 0, &val)
		v := float64(val)
		return v, strconv.FormatFloat(v, 'f', 0, 64), nil

	case TypeDWord:
		buf := make([]byte, 4)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		var val uint32
		helper.GetValueAt(buf, 0, &val)
		v := float64(val)
		return v, strconv.FormatFloat(v, 'f', 0, 64), nil

	case TypeInt:
		buf := make([]byte, 2)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		var val int16
		helper.GetValueAt(buf, 0, &val)
		v := float64(val)
		return v, strconv.FormatFloat(v, 'f', 0, 64), nil

	case TypeDInt:
		buf := make([]byte, 4)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		var val int32
		helper.GetValueAt(buf, 0, &val)
		v := float64(val)
		return v, strconv.FormatFloat(v, 'f', 0, 64), nil

	case TypeReal:
		buf := make([]byte, 4)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		val := helper.GetRealAt(buf, 0)
		v := float64(val)
		return v, strconv.FormatFloat(v, 'f', 6, 64), nil

	case TypeString:
		// S7 String: byte0=maxLen, byte1=actualLen, byte2+= chars
		buf := make([]byte, 256)
		if err := c.readArea(req, buf); err != nil {
			return 0, "", err
		}
		actualLen := int(buf[1])
		if actualLen > 254 {
			actualLen = 254
		}
		str := string(buf[2 : 2+actualLen])
		return 0, str, nil

	default:
		return 0, "", fmt.Errorf("unsupported data type: %s", req.DataType)
	}
}

// readArea 按存储区类型调用对应的 gos7 读取方法
func (c *Client) readArea(req ReadRequest, buf []byte) error {
	switch req.Area {
	case AreaDB:
		return c.s7client.AGReadDB(req.DBNumber, req.ByteOffset, len(buf), buf)
	case AreaM:
		return c.s7client.AGReadMB(req.ByteOffset, len(buf), buf)
	case AreaI:
		return c.s7client.AGReadEB(req.ByteOffset, len(buf), buf)
	case AreaQ:
		return c.s7client.AGReadAB(req.ByteOffset, len(buf), buf)
	case AreaV:
		// S7-200 SMART V 区等同于 DB1
		return c.s7client.AGReadDB(1, req.ByteOffset, len(buf), buf)
	default:
		return fmt.Errorf("unsupported area: %s", req.Area)
	}
}
