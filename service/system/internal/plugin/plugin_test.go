package plugin

import (
	"context"
	"testing"
)

// MockPlugin 模拟插件
type MockPlugin struct {
	name       string
	pluginType string
	platform   string
	started    bool
}

func (p *MockPlugin) Meta() Meta {
	return Meta{
		Name: p.name,
		Type: p.pluginType,
	}
}

func (p *MockPlugin) Init(ctx context.Context, config map[string]interface{}) error {
	return nil
}

func (p *MockPlugin) Start(ctx context.Context) error {
	p.started = true
	return nil
}

func (p *MockPlugin) Stop(ctx context.Context) error {
	p.started = false
	return nil
}

func (p *MockPlugin) HealthCheck(ctx context.Context) error {
	return nil
}

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()
	p := &MockPlugin{name: "test-plugin", pluginType: "channel"}

	if err := registry.Register(p); err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// 重复注册应该失败
	if err := registry.Register(p); err == nil {
		t.Error("expected error for duplicate registration")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	p := &MockPlugin{name: "test-plugin", pluginType: "channel"}
	registry.Register(p)

	got, err := registry.Get("test-plugin")
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if got.Meta().Name != "test-plugin" {
		t.Errorf("expected test-plugin, got %s", got.Meta().Name)
	}

	_, err = registry.Get("non-existent")
	if err == nil {
		t.Error("expected error for non-existent plugin")
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&MockPlugin{name: "p1", pluginType: "channel"})
	registry.Register(&MockPlugin{name: "p2", pluginType: "ai_model"})

	list := registry.List()
	if len(list) != 2 {
		t.Errorf("expected 2 plugins, got %d", len(list))
	}
}

func TestRegistry_GetByType(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&MockPlugin{name: "p1", pluginType: "channel"})
	registry.Register(&MockPlugin{name: "p2", pluginType: "ai_model"})
	registry.Register(&MockPlugin{name: "p3", pluginType: "channel"})

	channels := registry.GetByType("channel")
	if len(channels) != 2 {
		t.Errorf("expected 2 channel plugins, got %d", len(channels))
	}

	aiModels := registry.GetByType("ai_model")
	if len(aiModels) != 1 {
		t.Errorf("expected 1 ai_model plugin, got %d", len(aiModels))
	}
}

func TestRegistry_EnableDisable(t *testing.T) {
	registry := NewRegistry()
	p := &MockPlugin{name: "test-plugin", pluginType: "channel"}
	registry.Register(p)

	if err := registry.Enable("test-plugin"); err != nil {
		t.Fatalf("enable failed: %v", err)
	}

	info, _ := registry.GetPluginInfo("test-plugin")
	if info.Status != StatusEnabled {
		t.Errorf("expected enabled, got %s", info.Status)
	}

	if err := registry.Disable("test-plugin"); err != nil {
		t.Fatalf("disable failed: %v", err)
	}

	info, _ = registry.GetPluginInfo("test-plugin")
	if info.Status != StatusDisabled {
		t.Errorf("expected disabled, got %s", info.Status)
	}
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry()
	p := &MockPlugin{name: "test-plugin", pluginType: "channel"}
	registry.Register(p)

	if err := registry.Unregister("test-plugin"); err != nil {
		t.Fatalf("unregister failed: %v", err)
	}

	_, err := registry.Get("test-plugin")
	if err == nil {
		t.Error("expected error after unregister")
	}
}

func TestRegistry_RecordCall(t *testing.T) {
	registry := NewRegistry()
	p := &MockPlugin{name: "test-plugin", pluginType: "channel"}
	registry.Register(p)

	registry.RecordCall("test-plugin", 100.0)
	registry.RecordCall("test-plugin", 200.0)

	info, _ := registry.GetPluginInfo("test-plugin")
	if info.CallCount != 2 {
		t.Errorf("expected 2 calls, got %d", info.CallCount)
	}
}

func TestMarketplace_GetSDKTemplate(t *testing.T) {
	registry := NewRegistry()
	marketplace := NewMarketplace(registry)

	tests := []struct {
		pluginType string
		language   string
		wantErr    bool
	}{
		{"channel", "go", false},
		{"ai_model", "go", false},
		{"analyzer", "go", false},
		{"channel", "python", true},
		{"unknown", "go", true},
	}

	for _, tt := range tests {
		t.Run(tt.pluginType+"/"+tt.language, func(t *testing.T) {
			_, err := marketplace.GetSDKTemplate(context.Background(), tt.pluginType, tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestMarketplace_SDKTemplateContent(t *testing.T) {
	registry := NewRegistry()
	marketplace := NewMarketplace(registry)

	tmpl, err := marketplace.GetSDKTemplate(context.Background(), "channel", "go")
	if err != nil {
		t.Fatalf("get template failed: %v", err)
	}

	if tmpl.PluginType != "channel" {
		t.Errorf("expected channel, got %s", tmpl.PluginType)
	}
	if tmpl.Template == "" {
		t.Error("expected non-empty template")
	}
	if tmpl.Example == "" {
		t.Error("expected non-empty example")
	}
}

func TestRegisterBuiltinPlugins(t *testing.T) {
	registry := NewRegistry()

	if err := RegisterBuiltinPlugins(registry); err != nil {
		t.Fatalf("register builtin failed: %v", err)
	}

	// 检查内置插件数量
	list := registry.List()
	if len(list) < 10 {
		t.Errorf("expected at least 10 builtin plugins, got %d", len(list))
	}

	// 检查渠道插件
	channels := registry.GetByType("channel")
	if len(channels) < 5 {
		t.Errorf("expected at least 5 channel plugins, got %d", len(channels))
	}

	// 检查AI模型插件
	aiModels := registry.GetByType("ai_model")
	if len(aiModels) < 4 {
		t.Errorf("expected at least 4 ai_model plugins, got %d", len(aiModels))
	}

	// 检查分析器插件
	analyzers := registry.GetByType("analyzer")
	if len(analyzers) < 3 {
		t.Errorf("expected at least 3 analyzer plugins, got %d", len(analyzers))
	}
}

func TestValidatePluginConfig(t *testing.T) {
	schema := `{"required": ["api_key", "endpoint"]}`

	tests := []struct {
		config  map[string]interface{}
		wantErr bool
		desc    string
	}{
		{
			config:  map[string]interface{}{"api_key": "test", "endpoint": "https://api.example.com"},
			wantErr: false,
			desc:    "all required fields present",
		},
		{
			config:  map[string]interface{}{"api_key": "test"},
			wantErr: true,
			desc:    "missing required field",
		},
		{
			config:  map[string]interface{}{},
			wantErr: true,
			desc:    "empty config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := ValidatePluginConfig(schema, tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("error: %v, wantErr: %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginStatus(t *testing.T) {
	statuses := []PluginStatus{
		StatusRegistered,
		StatusEnabled,
		StatusDisabled,
		StatusError,
	}

	for _, s := range statuses {
		if s == "" {
			t.Error("empty status")
		}
	}
}

func TestMeta_Validation(t *testing.T) {
	meta := Meta{
		Name:        "test-plugin",
		Version:     "1.0.0",
		Type:        "channel",
		Author:      "Test Author",
		Description: "Test Description",
		License:     "MIT",
	}

	if meta.Name == "" {
		t.Error("expected non-empty name")
	}
	if meta.Type != "channel" {
		t.Errorf("expected channel, got %s", meta.Type)
	}
}
