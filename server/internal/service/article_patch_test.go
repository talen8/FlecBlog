package service

import (
	"testing"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
)

func stringPtr(v string) *string { return &v }
func uintPtr(v uint) *uint       { return &v }
func boolPtr(v bool) *bool       { return &v }

func TestApplyArticlePatchPreservesOmittedFields(t *testing.T) {
	categoryID := uint(8)
	article := &model.Article{
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
		CategoryID: &categoryID,
	}

	applyArticlePatch(article, &dto.PatchArticleRequest{
		Content: stringPtr("new content"),
	})

	if article.Content != "new content" {
		t.Fatalf("content = %q, want %q", article.Content, "new content")
	}
	if article.Summary != "old summary" || article.AISummary != "old ai summary" {
		t.Fatalf("summary fields changed unexpectedly: summary=%q ai_summary=%q", article.Summary, article.AISummary)
	}
	if article.Cover != "old.png" || article.Location != "Shenzhen" {
		t.Fatalf("metadata changed unexpectedly: cover=%q location=%q", article.Cover, article.Location)
	}
	if article.CategoryID == nil || *article.CategoryID != 8 {
		t.Fatalf("category changed unexpectedly: %v", article.CategoryID)
	}
	if article.IsPublish || !article.IsTop || !article.IsEssence || article.IsOutdated {
		t.Fatalf("status fields changed unexpectedly: %+v", article)
	}
}

func TestApplyArticlePatchSupportsExplicitClear(t *testing.T) {
	categoryID := uint(8)
	article := &model.Article{
		Summary:    "old summary",
		AISummary:  "old ai summary",
		Cover:      "old.png",
		Location:   "Shenzhen",
		CategoryID: &categoryID,
	}

	applyArticlePatch(article, &dto.PatchArticleRequest{
		Summary:    stringPtr(""),
		AISummary:  stringPtr(""),
		Cover:      stringPtr(""),
		Location:   stringPtr(""),
		CategoryID: uintPtr(0),
	})

	if article.Summary != "" || article.AISummary != "" || article.Cover != "" || article.Location != "" {
		t.Fatalf("explicit clear was not applied: %+v", article)
	}
	if article.CategoryID != nil {
		t.Fatalf("category_id = %v, want nil", article.CategoryID)
	}
}

func TestApplyArticlePatchSetsCategoryAndStatus(t *testing.T) {
	article := &model.Article{}

	applyArticlePatch(article, &dto.PatchArticleRequest{
		CategoryID: uintPtr(12),
		IsPublish:  boolPtr(true),
		IsTop:      boolPtr(true),
		IsEssence:  boolPtr(true),
		IsOutdated: boolPtr(true),
	})

	if article.CategoryID == nil || *article.CategoryID != 12 {
		t.Fatalf("category_id = %v, want 12", article.CategoryID)
	}
	if !article.IsPublish || !article.IsTop || !article.IsEssence || !article.IsOutdated {
		t.Fatalf("status fields not applied: %+v", article)
	}
}
