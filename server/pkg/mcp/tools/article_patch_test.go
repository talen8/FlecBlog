package tools

import (
	"encoding/json"
	"testing"
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
