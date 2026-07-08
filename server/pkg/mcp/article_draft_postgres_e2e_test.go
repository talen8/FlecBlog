package mcp

import (
	"context"
	"testing"

	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/internal/testutil"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestGranularArticleDraftToolsOverMCPPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	category := model.Category{Name: "MCP Draft E2E Category", Slug: "mcp-draft-e2e-category"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	tags := []model.Tag{
		{Name: "MCP Draft E2E Tag A", Slug: "mcp-draft-e2e-tag-a"},
		{Name: "MCP Draft E2E Tag B", Slug: "mcp-draft-e2e-tag-b"},
	}
	if err := db.Create(&tags).Error; err != nil {
		t.Fatalf("create tags: %v", err)
	}

	published := model.Article{
		Title:      "MCP Draft E2E Published",
		Slug:       "mcp-draft-e2e-published",
		Content:    "published body",
		Summary:    "published summary",
		IsPublish:  true,
		CategoryID: &category.ID,
	}
	if err := db.Create(&published).Error; err != nil {
		t.Fatalf("create published article: %v", err)
	}

	articleRepo := repository.NewArticleRepository(db)
	tagRepo := repository.NewTagRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	articleService := service.NewArticleService(articleRepo, tagRepo, categoryRepo, commentRepo, nil, db)

	handler := NewPublicHandler(articleService, nil, nil, nil, nil, nil, nil, nil, nil)
	server := newMCPTestServer(t, handler)

	ctx := context.Background()
	client := sdkmcp.NewClient(&sdkmcp.Implementation{Name: "flecblog-mcp-test", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &sdkmcp.StreamableClientTransport{Endpoint: server.URL}, nil)
	if err != nil {
		t.Fatalf("connect MCP client: %v", err)
	}
	defer session.Close()

	createResult := callMCPTool(t, ctx, session, "article_create_draft", map[string]any{
		"title":       "MCP Draft E2E Created",
		"content":     "draft body",
		"summary":     "draft summary",
		"cover":       "draft-cover.png",
		"location":    "Shenzhen",
		"is_top":      true,
		"is_essence":  true,
		"is_outdated": true,
		"category_id": category.ID,
		"tag_ids":     []uint{tags[0].ID, tags[1].ID},
	})
	var createOutput tools.ArticleDraftOutput
	decodeStructuredContent(t, createResult.StructuredContent, &createOutput)
	if createOutput.Item.ID == 0 {
		t.Fatal("article_create_draft returned ID 0")
	}

	created := loadMCPArticle(t, db, createOutput.Item.ID)
	if created.IsPublish {
		t.Fatal("article_create_draft persisted is_publish=true, want false")
	}
	if created.Cover != "draft-cover.png" || created.Location != "Shenzhen" {
		t.Fatalf("created metadata = cover %q location %q", created.Cover, created.Location)
	}
	if !created.IsTop || !created.IsEssence || !created.IsOutdated {
		t.Fatalf("created flags = top %t essence %t outdated %t, want all true", created.IsTop, created.IsEssence, created.IsOutdated)
	}
	if created.CategoryID == nil || *created.CategoryID != category.ID {
		t.Fatalf("created category_id = %v, want %d", created.CategoryID, category.ID)
	}
	assertMCPTagIDs(t, created.Tags, []uint{tags[0].ID, tags[1].ID})

	updateResult := callMCPTool(t, ctx, session, "article_update_draft", map[string]any{
		"id":      created.ID,
		"summary": "updated draft summary",
	})
	var updateOutput tools.ArticleDraftOutput
	decodeStructuredContent(t, updateResult.StructuredContent, &updateOutput)
	if updateOutput.Item.IsPublish {
		t.Fatal("article_update_draft returned published article")
	}

	updated := loadMCPArticle(t, db, created.ID)
	if updated.Summary != "updated draft summary" {
		t.Fatalf("updated summary = %q, want %q", updated.Summary, "updated draft summary")
	}
	if updated.Cover != created.Cover || updated.Location != created.Location {
		t.Fatalf("omitted metadata changed: cover=%q location=%q", updated.Cover, updated.Location)
	}
	if !updated.IsTop || !updated.IsEssence || !updated.IsOutdated {
		t.Fatalf("omitted flags changed: top=%t essence=%t outdated=%t", updated.IsTop, updated.IsEssence, updated.IsOutdated)
	}
	if updated.CategoryID == nil || *updated.CategoryID != category.ID {
		t.Fatalf("omitted category_id = %v, want %d", updated.CategoryID, category.ID)
	}
	assertMCPTagIDs(t, updated.Tags, []uint{tags[0].ID, tags[1].ID})

	callMCPTool(t, ctx, session, "article_update_draft", map[string]any{
		"id":          created.ID,
		"is_top":      false,
		"is_essence":  false,
		"is_outdated": false,
	})

	flagsCleared := loadMCPArticle(t, db, created.ID)
	if flagsCleared.IsTop || flagsCleared.IsEssence || flagsCleared.IsOutdated {
		t.Fatalf("explicit false flags did not persist: top=%t essence=%t outdated=%t", flagsCleared.IsTop, flagsCleared.IsEssence, flagsCleared.IsOutdated)
	}

	callMCPTool(t, ctx, session, "article_update_draft", map[string]any{
		"id":          created.ID,
		"summary":     "",
		"ai_summary":  "",
		"cover":       "",
		"location":    "",
		"category_id": 0,
		"tag_ids":     []uint{},
	})

	cleared := loadMCPArticle(t, db, created.ID)
	if cleared.IsPublish {
		t.Fatal("article_update_draft explicit clear changed draft to published")
	}
	if cleared.Summary != "" || cleared.AISummary != "" || cleared.Cover != "" || cleared.Location != "" {
		t.Fatalf("explicit clears did not persist: summary=%q ai_summary=%q cover=%q location=%q", cleared.Summary, cleared.AISummary, cleared.Cover, cleared.Location)
	}
	if cleared.CategoryID != nil {
		t.Fatalf("category_id = %v, want nil after explicit clear", cleared.CategoryID)
	}
	if len(cleared.Tags) != 0 {
		t.Fatalf("tags = %v, want empty after explicit clear", cleared.Tags)
	}

	publishedUpdate, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name: "article_update_draft",
		Arguments: map[string]any{
			"id":      published.ID,
			"summary": "must not persist",
		},
	})
	if err != nil {
		t.Fatalf("call article_update_draft for published article: %v", err)
	}
	if !publishedUpdate.IsError {
		t.Fatal("article_update_draft published article result isError = false, want true")
	}

	var publishedAfter model.Article
	if err := db.First(&publishedAfter, published.ID).Error; err != nil {
		t.Fatalf("reload published article: %v", err)
	}
	if publishedAfter.Summary != published.Summary {
		t.Fatalf("published summary changed to %q, want %q", publishedAfter.Summary, published.Summary)
	}
}
