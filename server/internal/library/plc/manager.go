// Package plc 西门子 S7 PLC 通信驱动
package plc

import (
	"sync"
)

// Manager 多设备连接池，key = deviceID
type Manager struct {
	clients map[int]*Client
	mu      sync.RWMutex
}

var defaultManager = &Manager{
	clients: make(map[int]*Client),
}

// GetOrCreate 获取或创建指定设备的连接
func GetOrCreate(deviceID int, host string, port, rack, slot int) *Client {
	defaultManager.mu.RLock()
	c, ok := defaultManager.clients[deviceID]
	defaultManager.mu.RUnlock()
	if ok {
		return c
	}

	defaultManager.mu.Lock()
	defer defaultManager.mu.Unlock()
	// double-check
	if c, ok = defaultManager.clients[deviceID]; ok {
		return c
	}
	c = NewClient(host, port, rack, slot)
	defaultManager.clients[deviceID] = c
	return c
}

// Remove 移除并关闭指定设备连接（设备删除/禁用时调用）
func Remove(deviceID int) {
	defaultManager.mu.Lock()
	defer defaultManager.mu.Unlock()
	if c, ok := defaultManager.clients[deviceID]; ok {
		c.Close()
		delete(defaultManager.clients, deviceID)
	}
}

// CloseAll 关闭所有连接（进程退出时调用）
func CloseAll() {
	defaultManager.mu.Lock()
	defer defaultManager.mu.Unlock()
	for id, c := range defaultManager.clients {
		c.Close()
		delete(defaultManager.clients, id)
	}
}
