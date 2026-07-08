package tools

import (
	"encoding/json"
	"testing"
	"time"

	"flec_blog/internal/dto"
)

func TestArticleManageUpdatePreservesFieldPresence(t *testing.T) {
	var input ArticleManageInput
	if err := json.Unmarshal([]byte(`{"action":"update","payload":{"id":123,"is_publish":true}}`), &input); err != nil {
		t.Fatalf("unmarshal update payload: %v", err)
	}

	if input.Payload.Summary != nil {
		t.Fatalf("summary = %v, want nil when omitted", input.Payload.Summary)
	}
	if input.Payload.AISummary != nil {
		t.Fatalf("ai_summary = %v, want nil when omitted", input.Payload.AISummary)
	}
	if input.Payload.Cover != nil {
		t.Fatalf("cover = %v, want nil when omitted", input.Payload.Cover)
	}
	if input.Payload.Location != nil {
		t.Fatalf("location = %v, want nil when omitted", input.Payload.Location)
	}
	if input.Payload.CategoryID != nil {
		t.Fatalf("category_id = %v, want nil when omitted", input.Payload.CategoryID)
	}
	if input.Payload.TagIDs != nil {
		t.Fatalf("tag_ids = %v, want nil when omitted", input.Payload.TagIDs)
	}
}

func TestArticleManageUpdateDistinguishesExplicitEmptyValues(t *testing.T) {
	var input ArticleManageInput
	if err := json.Unmarshal([]byte(`{"action":"update","payload":{"id":123,"summary":"","ai_summary":"","cover":"","location":"","category_id":0,"tag_ids":[]}}`), &input); err != nil {
		t.Fatalf("unmarshal update payload: %v", err)
	}

	assertStringPointerValue(t, "summary", input.Payload.Summary, "")
	assertStringPointerValue(t, "ai_summary", input.Payload.AISummary, "")
	assertStringPointerValue(t, "cover", input.Payload.Cover, "")
	assertStringPointerValue(t, "location", input.Payload.Location, "")

	if input.Payload.CategoryID == nil || *input.Payload.CategoryID != 0 {
		t.Fatalf("category_id = %v, want non-nil 0", input.Payload.CategoryID)
	}
	if input.Payload.TagIDs == nil || len(input.Payload.TagIDs) != 0 {
		t.Fatalf("tag_ids = %v, want non-nil empty slice", input.Payload.TagIDs)
	}
}

func assertStringPointerValue(t *testing.T, name string, got *string, want string) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s = nil, want %q", name, want)
	}
	if *got != want {
		t.Fatalf("%s = %q, want %q", name, *got, want)
	}
}

func TestBuildArticleUpdateRequestPreservesOmittedFields(t *testing.T) {
	current := &dto.ArticleAdminDetailResponse{
		Title:      "old title",
		Content:    "old content",
		Summary:    "old summary",
		AISummary:  "old ai summary",
		Cover:      "old.png",
		Location:   "Shenzhen",
		IsPublish:  false,
		IsTop:      true,
		IsEssence:  true,
		IsOutdated: false,
	}
	current.Category.ID = 8

	publish := true
	req := buildArticleUpdateRequest(current, ArticleManagePayload{
		ID:        123,
		IsPublish: &publish,
	})

	if req.Title != "" || req.Content != "" {
		t.Fatalf("legacy title/content no-op contract changed: title=%q content=%q", req.Title, req.Content)
	}
	if req.Summary != "old summary" || req.AISummary != "old ai summary" {
		t.Fatalf("summary fields changed unexpectedly: summary=%q ai_summary=%q", req.Summary, req.AISummary)
	}
	if req.Cover != "old.png" || req.Location != "Shenzhen" {
		t.Fatalf("metadata changed unexpectedly: cover=%q location=%q", req.Cover, req.Location)
	}
	if req.CategoryID == nil || *req.CategoryID != 8 {
		t.Fatalf("category_id = %v, want 8", req.CategoryID)
	}
	if req.TagIDs != nil {
		t.Fatalf("tag_ids = %v, want nil when omitted", req.TagIDs)
	}
	if req.IsPublish == nil || !*req.IsPublish {
		t.Fatalf("is_publish = %v, want explicit true", req.IsPublish)
	}
}

func TestBuildArticleUpdateRequestSupportsExplicitClear(t *testing.T) {
	current := &dto.ArticleAdminDetailResponse{
		Summary:   "old summary",
		AISummary: "old ai summary",
		Cover:     "old.png",
		Location:  "Shenzhen",
	}
	current.Category.ID = 8

	empty := ""
	zero := uint(0)
	req := buildArticleUpdateRequest(current, ArticleManagePayload{
		ID:         123,
		Summary:    &empty,
		AISummary:  &empty,
		Cover:      &empty,
		Location:   &empty,
		CategoryID: &zero,
		TagIDs:     []uint{},
	})

	if req.Summary != "" || req.AISummary != "" || req.Cover != "" || req.Location != "" {
		t.Fatalf("explicit clear was not preserved: %+v", req)
	}
	if req.CategoryID != nil {
		t.Fatalf("category_id = %v, want nil for explicit 0", req.CategoryID)
	}
	if req.TagIDs == nil || len(req.TagIDs) != 0 {
		t.Fatalf("tag_ids = %v, want non-nil empty slice", req.TagIDs)
	}
}

func TestBuildArticleUpdateRequestSetsCategory(t *testing.T) {
	current := &dto.ArticleAdminDetailResponse{}
	categoryID := uint(12)

	req := buildArticleUpdateRequest(current, ArticleManagePayload{
		ID:         123,
		CategoryID: &categoryID,
	})

	if req.CategoryID == nil || *req.CategoryID != 12 {
		t.Fatalf("category_id = %v, want 12", req.CategoryID)
	}
}

func TestBuildArticleUpdateRequestCarriesInternalRevision(t *testing.T) {
	revision := time.Date(2026, time.July, 7, 22, 30, 0, 123000000, time.UTC)
	current := &dto.ArticleAdminDetailResponse{Revision: revision}

	req := buildArticleUpdateRequest(current, ArticleManagePayload{ID: 123})
	if req.ExpectedUpdatedAt == nil || !req.ExpectedUpdatedAt.Equal(revision) {
		t.Fatalf("expected_updated_at = %v, want %s", req.ExpectedUpdatedAt, revision)
	}

	encoded, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal request: %v", err)
	}
	if string(encoded) == "" || json.Valid(encoded) == false {
		t.Fatalf("invalid JSON: %q", encoded)
	}
	if containsJSONKey(encoded, "ExpectedUpdatedAt") || containsJSONKey(encoded, "expected_updated_at") {
		t.Fatalf("internal revision leaked into JSON: %s", encoded)
	}
}

func containsJSONKey(encoded []byte, key string) bool {
	var value map[string]any
	if err := json.Unmarshal(encoded, &value); err != nil {
		return false
	}
	_, ok := value[key]
	return ok
}
