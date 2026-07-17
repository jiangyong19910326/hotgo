// Package mqttx MQTT 订阅客户端
package mqttx

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	mqtt.ERROR = log.New(os.Stderr, "[paho-ERR] ", log.LstdFlags)
	mqtt.CRITICAL = log.New(os.Stderr, "[paho-CRT] ", log.LstdFlags)
	mqtt.WARN = log.New(os.Stderr, "[paho-WRN] ", log.LstdFlags)
}

type config struct {
	Addr     string `json:"addr"`
	ClientId string `json:"clientId"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
	QoS      byte   `json:"qos"`
}

var (
	client    mqtt.Client
	clientMu  sync.Mutex
	startOnce sync.Once
)

// Start 连接 broker, 加载缓存, 订阅 topic. 失败不阻塞主进程, 由 paho 自动重连。
func Start(ctx context.Context) {
	startOnce.Do(func() {
		go run(ctx)
	})
}

func run(ctx context.Context) {
	cfg, err := loadConfig(ctx)
	if err != nil {
		g.Log().Warningf(ctx, "mqttx disabled: %v", err)
		return
	}
	if cfg.Addr == "" {
		g.Log().Info(ctx, "mqttx disabled: empty addr")
		return
	}

	if err = Reload(ctx); err != nil {
		g.Log().Warningf(ctx, "mqttx initial cache load err: %v", err)
	}

	g.Log().Infof(ctx, "mqttx connecting: addr=%s clientId=%s topic=%s", cfg.Addr, cfg.ClientId, cfg.Topic)

	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Addr).
		SetClientID(cfg.ClientId).
		SetUsername(cfg.Username).
		SetPassword(cfg.Password).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetConnectTimeout(5 * time.Second).
		SetKeepAlive(30 * time.Second).
		SetCleanSession(true).
		SetOnConnectHandler(func(c mqtt.Client) {
			g.Log().Infof(ctx, "mqttx connected, subscribing topic: %s", cfg.Topic)
			if token := c.Subscribe(cfg.Topic, cfg.QoS, onMessage); token.Wait() && token.Error() != nil {
				g.Log().Errorf(ctx, "mqttx subscribe err: %v", token.Error())
			} else {
				g.Log().Infof(ctx, "mqttx subscribed ok: %s", cfg.Topic)
			}
		}).
		SetConnectionLostHandler(func(_ mqtt.Client, e error) {
			g.Log().Warningf(ctx, "mqttx connection lost: %v", e)
		}).
		SetReconnectingHandler(func(_ mqtt.Client, _ *mqtt.ClientOptions) {
			g.Log().Warning(ctx, "mqttx reconnecting...")
		}).
		SetDefaultPublishHandler(func(_ mqtt.Client, m mqtt.Message) {
			g.Log().Debugf(ctx, "mqttx msg(default): topic=%s len=%d", m.Topic(), len(m.Payload()))
		})

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.WaitTimeout(10*time.Second) && token.Error() != nil {
		g.Log().Errorf(ctx, "mqttx connect err: %v (will auto-retry)", token.Error())
	}

	clientMu.Lock()
	client = c
	clientMu.Unlock()

	go reloadLoop(ctx)
}

// reloadLoop 30秒刷一次设备/点位缓存, 避免管理后台改动后内存陈旧
func reloadLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := Reload(ctx); err != nil {
				g.Log().Warningf(ctx, "mqttx cache reload err: %v", err)
			}
		}
	}
}

// Stop 断开连接
func Stop() {
	clientMu.Lock()
	defer clientMu.Unlock()
	if client != nil && client.IsConnected() {
		client.Disconnect(500)
	}
}

// Publish 向 MQTT broker 发布消息，供设备控制命令复用当前连接。
func Publish(ctx context.Context, topic string, payload []byte, qos byte) error {
	clientMu.Lock()
	c := client
	clientMu.Unlock()

	if c == nil || !c.IsConnected() {
		g.Log().Warningf(ctx, "mqttx publish skipped: client not connected topic=%s payload=%s", topic, string(payload))
		return gerror.New("mqtt client not connected")
	}

	token := c.Publish(topic, qos, false, payload)
	if !token.WaitTimeout(5 * time.Second) {
		g.Log().Warningf(ctx, "mqttx publish timeout: topic=%s payload=%s", topic, string(payload))
		return gerror.New("mqtt publish timeout")
	}
	if token.Error() != nil {
		g.Log().Warningf(ctx, "mqttx publish failed: topic=%s err=%v payload=%s", topic, token.Error(), string(payload))
		return token.Error()
	}
	g.Log().Infof(ctx, "mqttx published: topic=%s payload=%s", topic, string(payload))
	return nil
}

func loadConfig(ctx context.Context) (*config, error) {
	v, err := g.Cfg().Get(ctx, "mqtt")
	if err != nil {
		return nil, err
	}
	cfg := &config{
		ClientId: "hotgo-mqtt-sub",
		QoS:      0,
	}
	if err = v.Scan(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
