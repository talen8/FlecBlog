package mcp

import (
	"sort"
	"testing"

	"flec_blog/internal/model"

	"gorm.io/gorm"
)

func loadMCPArticle(t *testing.T, db *gorm.DB, id uint) model.Article {
	t.Helper()
	var article model.Article
	if err := db.Preload("Category").Preload("Tags").First(&article, id).Error; err != nil {
		t.Fatalf("load article %d: %v", id, err)
	}
	return article
}

func assertMCPTagIDs(t *testing.T, tags []model.Tag, want []uint) {
	t.Helper()
	got := make([]uint, len(tags))
	for i, tag := range tags {
		got[i] = tag.ID
	}
	sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
	sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
	if len(got) != len(want) {
		t.Fatalf("tag IDs = %v, want %v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("tag IDs = %v, want %v", got, want)
		}
	}
}
