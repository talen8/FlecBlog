package service

import (
	"testing"

	"flec_blog/config"
	"flec_blog/internal/model"
	"flec_blog/internal/testutil"
)

func TestMCPStaticOperatorSettingOverPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	settingService := NewSettingService(db)
	cfg := &config.Config{}
	settingService.SetConfig(cfg)

	rotatedSecret, err := settingService.ResetMCPSecret()
	if err != nil {
		t.Fatalf("reset MCP secret: %v", err)
	}
	if rotatedSecret == "" {
		t.Fatal("reset MCP secret returned empty value")
	}
	if settingService.MCPSecret() != rotatedSecret {
		t.Fatal("MCP secret provider did not reflect reset")
	}
	secretSettings, err := settingService.GetByGroup(model.SettingGroupAI)
	if err != nil {
		t.Fatalf("get AI settings after secret reset: %v", err)
	}
	if secretSettings[KeyAIMCPSecret] != rotatedSecret {
		t.Fatal("persisted MCP secret did not match reset value")
	}

	settings, err := settingService.GetByGroup(model.SettingGroupAI)
	if err != nil {
		t.Fatalf("get AI settings: %v", err)
	}
	if _, ok := settings[KeyAIMCPStaticOperatorUserID]; !ok {
		t.Fatalf("migration did not create %s", KeyAIMCPStaticOperatorUserID)
	}

	if settings[KeyAIMCPAdminToolsEnabled] != "false" {
		t.Fatalf("migration default %s = %q, want false", KeyAIMCPAdminToolsEnabled, settings[KeyAIMCPAdminToolsEnabled])
	}
	if settingService.MCPAdminToolsEnabled() {
		t.Fatal("MCP admin tools unexpectedly enabled by default")
	}

	if err := settingService.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPAdminToolsEnabled: "true",
	}); err != nil {
		t.Fatalf("enable MCP admin tools: %v", err)
	}
	if !settingService.MCPAdminToolsEnabled() {
		t.Fatal("MCP admin tools provider did not hot-reload true")
	}

	if err := settingService.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPAdminToolsEnabled: "invalid",
	}); err == nil {
		t.Fatal("invalid MCP admin tools value error = nil")
	}
	settings, err = settingService.GetByGroup(model.SettingGroupAI)
	if err != nil {
		t.Fatalf("reload AI settings after invalid admin-tools update: %v", err)
	}
	if settings[KeyAIMCPAdminToolsEnabled] != "true" {
		t.Fatalf("invalid admin-tools update persisted value %q, want true", settings[KeyAIMCPAdminToolsEnabled])
	}

	if err := settingService.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPAdminToolsEnabled: "false",
	}); err != nil {
		t.Fatalf("disable MCP admin tools: %v", err)
	}
	if settingService.MCPAdminToolsEnabled() {
		t.Fatal("MCP admin tools provider did not hot-reload false")
	}

	if err := settingService.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPStaticOperatorUserID: "42",
	}); err != nil {
		t.Fatalf("set operator user ID: %v", err)
	}
	if cfg.AI.MCPStaticOperatorUserID != 42 {
		t.Fatalf("hot-reloaded operator user ID = %d, want 42", cfg.AI.MCPStaticOperatorUserID)
	}

	if err := settingService.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPStaticOperatorUserID: "not-an-id",
	}); err == nil {
		t.Fatal("invalid operator user ID error = nil")
	}
	settings, err = settingService.GetByGroup(model.SettingGroupAI)
	if err != nil {
		t.Fatalf("reload AI settings after invalid update: %v", err)
	}
	if settings[KeyAIMCPStaticOperatorUserID] != "42" {
		t.Fatalf("invalid update persisted value %q, want 42", settings[KeyAIMCPStaticOperatorUserID])
	}

	if err := settingService.UpdateGroup(model.SettingGroupAI, map[string]string{
		KeyAIMCPStaticOperatorUserID: "",
	}); err != nil {
		t.Fatalf("reset operator user ID to auto: %v", err)
	}
	if cfg.AI.MCPStaticOperatorUserID != 0 {
		t.Fatalf("hot-reloaded auto operator user ID = %d, want 0", cfg.AI.MCPStaticOperatorUserID)
	}
}
