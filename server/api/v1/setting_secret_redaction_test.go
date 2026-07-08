package v1

import (
	"testing"

	"flec_blog/internal/model"
	"flec_blog/internal/service"

	"github.com/gin-gonic/gin"
)

func TestRedactMCPSecretForCaller(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		group      string
		user       *model.User
		wantSecret bool
	}{
		{name: "ordinary admin ai group", group: model.SettingGroupAI, user: &model.User{Role: model.RoleAdmin}, wantSecret: false},
		{name: "super admin ai group", group: model.SettingGroupAI, user: &model.User{Role: model.RoleSuperAdmin}, wantSecret: true},
		{name: "missing user fails closed", group: model.SettingGroupAI, user: nil, wantSecret: false},
		{name: "other group unchanged", group: model.SettingGroupBasic, user: &model.User{Role: model.RoleAdmin}, wantSecret: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx, _ := gin.CreateTestContext(nil)
			if tc.user != nil {
				ctx.Set("user", tc.user)
			}
			settings := map[string]string{
				service.KeyAIMCPSecret: "top-secret",
				service.KeyAIModel:     "model-1",
			}

			got := redactMCPSecretForCaller(ctx, tc.group, settings)
			_, hasSecret := got[service.KeyAIMCPSecret]
			if hasSecret != tc.wantSecret {
				t.Fatalf("secret present = %v, want %v; settings=%v", hasSecret, tc.wantSecret, got)
			}
			if got[service.KeyAIModel] != "model-1" {
				t.Fatalf("non-secret AI setting changed: %v", got)
			}
		})
	}
}

func TestRedactMCPSecretDoesNotMutateSourceMap(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(nil)
	ctx.Set("user", &model.User{Role: model.RoleAdmin})
	source := map[string]string{service.KeyAIMCPSecret: "top-secret"}

	got := redactMCPSecretForCaller(ctx, model.SettingGroupAI, source)
	if _, ok := got[service.KeyAIMCPSecret]; ok {
		t.Fatalf("redacted result leaked secret: %v", got)
	}
	if source[service.KeyAIMCPSecret] != "top-secret" {
		t.Fatalf("source map mutated: %v", source)
	}
}
