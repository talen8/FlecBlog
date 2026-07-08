package mcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"flec_blog/config"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/internal/testutil"
	mcpauth "flec_blog/pkg/mcp/auth"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGranularUserToolsOverMCPPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)
	t.Setenv("MCP_STATIC_OPERATOR_USER_ID", "")

	// --- Phase 0: Observe the migration baseline without mutating unrelated users. ---
	var activeSuperAdmins []model.User
	if err := db.Where(
		"role = ? AND is_enabled = ? AND deleted_at IS NULL",
		model.RoleSuperAdmin,
		true,
	).Find(&activeSuperAdmins).Error; err != nil {
		t.Fatalf("query active super admins: %v", err)
	}
	if len(activeSuperAdmins) != 1 {
		t.Fatalf("migration baseline active super_admin count = %d, want 1", len(activeSuperAdmins))
	}
	migrationSuperAdmin := activeSuperAdmins[0]
	t.Logf("migration-seeded super_admin: id=%d email=%s", migrationSuperAdmin.ID, migrationSuperAdmin.Email)

	// --- Seed non-super-admin test users (admin + regular target) ---

	admin := model.User{
		Email:       "mcp-admin@example.test",
		Nickname:    "MCP Admin",
		Role:        model.RoleAdmin,
		IsEnabled:   true,
		HasPassword: true,
	}
	target := model.User{
		Email:       "mcp-target@example.test",
		Nickname:    "MCP Target",
		Role:        model.RoleUser,
		IsEnabled:   true,
		HasPassword: true,
	}
	for _, user := range []*model.User{&admin, &target} {
		if err := db.Create(user).Error; err != nil {
			t.Fatalf("create seed user %s: %v", user.Email, err)
		}
	}
	t.Cleanup(func() {
		_ = db.Unscoped().Delete(&model.User{}, []uint{admin.ID, target.ID}).Error
	})

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, nil, &config.Config{JWT: config.JWTConfig{Secret: "mcp-user-e2e-jwt"}})
	handler := NewPublicHandlerWithOptions(
		nil, nil, nil, nil, nil, nil, nil, userService, nil,
		PublicHandlerOptions{AdminToolsEnabledProvider: func() bool { return true }},
	)

	// --- Phase 1: Single super_admin static operator (no env var needed) ---
	staticServer := newMCPTestServer(t, handler)
	staticSession := connectMCPTestClient(t, staticServer.URL)

	listResult, err := staticSession.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      "user_list",
		Arguments: map[string]any{},
	})
	if err != nil {
		t.Fatalf("call user_list: %v", err)
	}
	if listResult.IsError {
		content, _ := json.Marshal(listResult.Content)
		t.Fatalf("user_list returned isError=true: %s", content)
	}
	var listOutput tools.UserListOutput
	decodeStructuredContent(t, listResult.StructuredContent, &listOutput)
	if listOutput.Page != 1 || listOutput.PageSize != 20 {
		t.Fatalf("user_list default pagination = page %d size %d, want page 1 size %d", listOutput.Page, listOutput.PageSize, 20)
	}
	if !containsUserID(listOutput.Items, migrationSuperAdmin.ID) || !containsUserID(listOutput.Items, target.ID) {
		t.Fatalf("user_list IDs = %v", userIDs(listOutput.Items))
	}

	getResult := callMCPTool(t, context.Background(), staticSession, "user_get", map[string]any{"id": target.ID})
	var getOutput tools.UserGetOutput
	decodeStructuredContent(t, getResult.StructuredContent, &getOutput)
	if getOutput.Item.ID != target.ID || !getOutput.Item.IsEnabled {
		t.Fatalf("user_get item = %+v", getOutput.Item)
	}

	createResult := callMCPTool(t, context.Background(), staticSession, "user_create", map[string]any{
		"email":    "mcp-created@example.test",
		"password": "secret123",
		"nickname": "MCP Created",
		"role":     "user",
	})
	var createOutput tools.UserCreateOutput
	decodeStructuredContent(t, createResult.StructuredContent, &createOutput)
	if createOutput.Item.ID == 0 || createOutput.Item.Role != string(model.RoleUser) || !createOutput.Item.IsEnabled || !createOutput.Item.HasPassword {
		t.Fatalf("user_create item = %+v", createOutput.Item)
	}
	createdID := createOutput.Item.ID
	t.Cleanup(func() {
		_ = db.Unscoped().Delete(&model.User{}, createdID).Error
	})

	if err := db.Model(&model.User{}).Where("id = ?", createdID).Update("has_password", false).Error; err != nil {
		t.Fatalf("clear has_password before password-update regression: %v", err)
	}

	updateResult := callMCPTool(t, context.Background(), staticSession, "user_update", map[string]any{
		"id":         createdID,
		"nickname":   "MCP Updated",
		"role":       "guest",
		"is_enabled": false,
		"password":   "secret456",
	})
	var updateOutput tools.UserUpdateOutput
	decodeStructuredContent(t, updateResult.StructuredContent, &updateOutput)
	if updateOutput.Item.Nickname != "MCP Updated" || updateOutput.Item.Role != string(model.RoleGuest) || updateOutput.Item.IsEnabled || !updateOutput.Item.HasPassword {
		t.Fatalf("user_update item = %+v", updateOutput.Item)
	}

	deleteResult := callMCPTool(t, context.Background(), staticSession, "user_delete", map[string]any{"id": target.ID})
	var deleteOutput tools.UserDeleteOutput
	decodeStructuredContent(t, deleteResult.StructuredContent, &deleteOutput)
	if !deleteOutput.Deleted || deleteOutput.ID != target.ID {
		t.Fatalf("user_delete output = %+v", deleteOutput)
	}
	var deletedTarget model.User
	if err := db.Unscoped().First(&deletedTarget, target.ID).Error; err != nil {
		t.Fatalf("reload deleted target: %v", err)
	}
	if !deletedTarget.DeletedAt.Valid {
		t.Fatal("user_delete did not soft-delete target")
	}

	// --- Phase 2: Create second super_admin → resolver must fail closed ---
	secondSuperAdmin := model.User{
		Email:     "mcp-superadmin-2@example.test",
		Nickname:  "MCP Super Admin 2",
		Role:      model.RoleSuperAdmin,
		IsEnabled: true,
	}
	if err := db.Create(&secondSuperAdmin).Error; err != nil {
		t.Fatalf("create second super admin: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Unscoped().Delete(&model.User{}, secondSuperAdmin.ID).Error
	})
	ambiguous, err := staticSession.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name:      "user_list",
		Arguments: map[string]any{"page_size": 10},
	})
	if err != nil {
		t.Fatalf("call user_list with ambiguous static operator: %v", err)
	}
	if !ambiguous.IsError {
		t.Fatal("ambiguous static operator result isError = false, want true")
	}

	// --- Phase 3: WebAdmin-configured operator resolves the same ambiguity. ---
	configuredOperatorID := migrationSuperAdmin.ID
	providerHandler := NewPublicHandlerWithOptions(
		nil, nil, nil, nil, nil, nil, nil, userService, nil,
		PublicHandlerOptions{
			StaticOperatorUserIDProvider: func() uint { return configuredOperatorID },
			AdminToolsEnabledProvider:    func() bool { return true },
		},
	)
	providerServer := newMCPTestServer(t, providerHandler)
	providerSession := connectMCPTestClient(t, providerServer.URL)
	callMCPTool(t, context.Background(), providerSession, "user_list", map[string]any{"page_size": 10})

	// --- Phase 4: ENV override takes precedence over the WebAdmin provider. ---
	overrideEmail := "mcp-env-override-admin@example.test"
	t.Cleanup(func() {
		_ = db.Unscoped().Where("email = ?", overrideEmail).Delete(&model.User{}).Error
	})
	t.Setenv("MCP_STATIC_OPERATOR_USER_ID", strconv.FormatUint(uint64(admin.ID), 10))
	envOverrideResult, err := providerSession.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name: "user_create",
		Arguments: map[string]any{
			"email":    overrideEmail,
			"password": "secret123",
			"nickname": "ENV Override Admin Attempt",
			"role":     "admin",
		},
	})
	if err != nil {
		t.Fatalf("call user_create with ENV override: %v", err)
	}
	if !envOverrideResult.IsError {
		t.Fatal("ENV override did not take precedence over WebAdmin super_admin provider")
	}

	// --- Phase 5: Explicit ENV super_admin ID resolves ambiguity on the legacy path. ---
	t.Setenv("MCP_STATIC_OPERATOR_USER_ID", strconv.FormatUint(uint64(migrationSuperAdmin.ID), 10))
	callMCPTool(t, context.Background(), staticSession, "user_list", map[string]any{"page_size": 10})

	// --- Phase 6: OAuth operator scenarios ---
	oauthPrincipal := &mcpauth.Principal{
		Method:  "oauth",
		Subject: "oauth-admin-subject",
		Scopes:  map[string]struct{}{mcpauth.ScopeAdmin: {}},
		Claims:  map[string]any{"mcp_user_id": float64(admin.ID)},
	}
	if err := oauthPrincipal.BindVerifiedLocalUser(admin.ID); err != nil {
		t.Fatalf("bind verified OAuth operator: %v", err)
	}
	oauthServer := newMCPPrincipalTestServer(t, handler, oauthPrincipal)
	oauthSession := connectMCPTestClient(t, oauthServer.URL)

	allowedCreate := callMCPTool(t, context.Background(), oauthSession, "user_create", map[string]any{
		"email":    "mcp-oauth-created@example.test",
		"password": "secret123",
		"nickname": "OAuth Created",
		"role":     "user",
	})
	var oauthCreateOutput tools.UserCreateOutput
	decodeStructuredContent(t, allowedCreate.StructuredContent, &oauthCreateOutput)
	if oauthCreateOutput.Item.ID == 0 {
		t.Fatal("OAuth admin failed to create ordinary user")
	}
	t.Cleanup(func() {
		_ = db.Unscoped().Delete(&model.User{}, oauthCreateOutput.Item.ID).Error
	})

	forbiddenCreate, err := oauthSession.CallTool(context.Background(), &sdkmcp.CallToolParams{
		Name: "user_create",
		Arguments: map[string]any{
			"email":    "mcp-oauth-admin@example.test",
			"password": "secret123",
			"nickname": "OAuth Admin Attempt",
			"role":     "admin",
		},
	})
	if err != nil {
		t.Fatalf("call OAuth admin create-admin attempt: %v", err)
	}
	if !forbiddenCreate.IsError {
		t.Fatal("real local admin created another admin, want rejection")
	}
}

func newMCPPrincipalTestServer(t *testing.T, handler http.Handler, principal *mcpauth.Principal) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := mcpauth.ContextWithPrincipal(r.Context(), principal)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}))
	t.Cleanup(server.Close)
	return server
}

func connectMCPTestClient(t *testing.T, endpoint string) *sdkmcp.ClientSession {
	t.Helper()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-user-e2e", Version: "1.0"}, nil)
	session, err := client.Connect(context.Background(), &sdkmcp.StreamableClientTransport{Endpoint: endpoint}, nil)
	if err != nil {
		t.Fatalf("connect MCP test client: %v", err)
	}
	t.Cleanup(func() { _ = session.Close() })
	return session
}

func containsUserID(items []tools.UserItem, id uint) bool {
	for _, item := range items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func userIDs(items []tools.UserItem) []uint {
	ids := make([]uint, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	return ids
}
