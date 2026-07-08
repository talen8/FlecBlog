package tools

import (
	"context"
	"testing"
)

func TestNormalizeArticleReadPage(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "defaults", page: 0, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{name: "keeps bounded values", page: 3, pageSize: 50, wantPage: 3, wantPageSize: 50},
		{name: "caps page size", page: 2, pageSize: 500, wantPage: 2, wantPageSize: 100},
		{name: "negative values use defaults", page: -1, pageSize: -1, wantPage: 1, wantPageSize: 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageSize := normalizeArticleReadPage(tt.page, tt.pageSize)
			if gotPage != tt.wantPage || gotPageSize != tt.wantPageSize {
				t.Fatalf("normalizeArticleReadPage(%d, %d) = (%d, %d), want (%d, %d)",
					tt.page, tt.pageSize, gotPage, gotPageSize, tt.wantPage, tt.wantPageSize)
			}
		})
	}
}

func TestSearchArticleToolRejectsBlankKeyword(t *testing.T) {
	wrapper := NewArticleWrapper(nil)
	_, _, err := wrapper.SearchArticleTool(context.Background(), nil, ArticleSearchInput{Keyword: "   "})
	if err == nil {
		t.Fatal("SearchArticleTool() error = nil, want non-nil for blank keyword")
	}
}
