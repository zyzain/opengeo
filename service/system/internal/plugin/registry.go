package plugin

import (
	"fmt"
	"sync"
	"time"
)

// DefaultRegistry 默认插件注册表实现
type DefaultRegistry struct {
	mu       sync.RWMutex
	plugins  map[string]Plugin
	statuses map[string]PluginStatus
	infos    map[string]*PluginInfo
}

// NewRegistry 创建插件注册表
func NewRegistry() *DefaultRegistry {
	return &DefaultRegistry{
		plugins:  make(map[string]Plugin),
		statuses: make(map[string]PluginStatus),
		infos:    make(map[string]*PluginInfo),
	}
}

// Register 注册插件
func (r *DefaultRegistry) Register(p Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	meta := p.Meta()
	if meta.Name == "" {
		return fmt.Errorf("plugin name is empty")
	}

	if _, exists := r.plugins[meta.Name]; exists {
		return fmt.Errorf("plugin %s already registered", meta.Name)
	}

	r.plugins[meta.Name] = p
	r.statuses[meta.Name] = StatusRegistered
	r.infos[meta.Name] = &PluginInfo{
		Meta:     meta,
		Status:   StatusRegistered,
		LoadTime: time.Now().Unix(),
	}

	return nil
}

// Unregister 注销插件
func (r *DefaultRegistry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	delete(r.plugins, name)
	delete(r.statuses, name)
	delete(r.infos, name)

	return nil
}

// Get 获取插件
func (r *DefaultRegistry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return p, nil
}

// List 列出所有插件
func (r *DefaultRegistry) List() []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		result = append(result, p)
	}
	return result
}

// GetByType 按类型获取插件
func (r *DefaultRegistry) GetByType(pluginType string) []Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]Plugin, 0)
	for _, p := range r.plugins {
		if p.Meta().Type == pluginType {
			result = append(result, p)
		}
	}
	return result
}

// Enable 启用插件
func (r *DefaultRegistry) Enable(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	r.statuses[name] = StatusEnabled
	if info, ok := r.infos[name]; ok {
		info.Status = StatusEnabled
	}

	return nil
}

// Disable 禁用插件
func (r *DefaultRegistry) Disable(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plugins[name]; !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	r.statuses[name] = StatusDisabled
	if info, ok := r.infos[name]; ok {
		info.Status = StatusDisabled
	}

	return nil
}

// GetPluginInfo 获取插件运行时信息
func (r *DefaultRegistry) GetPluginInfo(name string) (*PluginInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, exists := r.infos[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return info, nil
}

// ListPluginInfos 列出所有插件运行时信息
func (r *DefaultRegistry) ListPluginInfos() []*PluginInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*PluginInfo, 0, len(r.infos))
	for _, info := range r.infos {
		result = append(result, info)
	}
	return result
}

// RecordCall 记录插件调用
func (r *DefaultRegistry) RecordCall(name string, latencyMs float64) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if info, ok := r.infos[name]; ok {
		info.CallCount++
		// 计算移动平均延迟
		if info.AvgLatency == 0 {
			info.AvgLatency = latencyMs
		} else {
			info.AvgLatency = (info.AvgLatency*0.9 + latencyMs*0.1)
		}
	}
}

// SetError 设置插件错误状态
func (r *DefaultRegistry) SetError(name string, err string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.statuses[name] = StatusError
	if info, ok := r.infos[name]; ok {
		info.Status = StatusError
		info.Error = err
	}
}
