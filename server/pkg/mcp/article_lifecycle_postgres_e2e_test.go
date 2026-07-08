package mcp

import (
	"context"
	"errors"
	"testing"

	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/internal/testutil"
	"flec_blog/pkg/mcp/tools"

	sdkmcp "github.com/modelcontextprotocol/go-sdk/mcp"
	"gorm.io/gorm"
)

func TestArticleLifecycleToolsOverMCPPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	category := model.Category{Name: "MCP Lifecycle E2E Category", Slug: "mcp-lifecycle-e2e-category"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	tags := []model.Tag{
		{Name: "MCP Lifecycle E2E Tag A", Slug: "mcp-lifecycle-e2e-tag-a"},
		{Name: "MCP Lifecycle E2E Tag B", Slug: "mcp-lifecycle-e2e-tag-b"},
	}
	if err := db.Create(&tags).Error; err != nil {
		t.Fatalf("create tags: %v", err)
	}

	published := model.Article{
		Title:      "MCP Lifecycle E2E Published",
		Slug:       "mcp-lifecycle-e2e-published",
		Content:    "published lifecycle body",
		Summary:    "published lifecycle summary",
		AISummary:  "published lifecycle ai summary",
		Cover:      "lifecycle-cover.png",
		Location:   "Shenzhen",
		IsPublish:  true,
		CategoryID: &category.ID,
	}
	if err := db.Create(&published).Error; err != nil {
		t.Fatalf("create published article: %v", err)
	}
	if err := db.Model(&published).Association("Tags").Replace(tags); err != nil {
		t.Fatalf("attach published tags: %v", err)
	}

	draft := model.Article{
		Title:     "MCP Lifecycle E2E Draft",
		Slug:      "mcp-lifecycle-e2e-draft",
		Content:   "draft lifecycle body",
		Summary:   "draft lifecycle summary",
		IsPublish: false,
	}
	if err := db.Create(&draft).Error; err != nil {
		t.Fatalf("create draft article: %v", err)
	}

	deletable := model.Article{
		Title:     "MCP Lifecycle E2E Delete",
		Slug:      "mcp-lifecycle-e2e-delete",
		Content:   "delete lifecycle body",
		IsPublish: false,
	}
	if err := db.Create(&deletable).Error; err != nil {
		t.Fatalf("create deletable article: %v", err)
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

	updateResult := callMCPTool(t, ctx, session, "article_update_published", map[string]any{
		"id":          published.ID,
		"summary":     "updated published summary",
		"is_top":      true,
		"is_essence":  true,
		"is_outdated": true,
	})
	var updateOutput tools.ArticleUpdatePublishedOutput
	decodeStructuredContent(t, updateResult.StructuredContent, &updateOutput)
	if !updateOutput.Item.IsPublish {
		t.Fatal("article_update_published returned draft")
	}

	updated := loadMCPArticle(t, db, published.ID)
	if !updated.IsPublish {
		t.Fatal("article_update_published persisted is_publish=false")
	}
	if updated.Summary != "updated published summary" {
		t.Fatalf("updated summary = %q", updated.Summary)
	}
	if updated.Title != published.Title || updated.Content != published.Content {
		t.Fatalf("omitted title/content changed: title=%q content=%q", updated.Title, updated.Content)
	}
	if updated.Cover != published.Cover || updated.Location != published.Location {
		t.Fatalf("omitted metadata changed: cover=%q location=%q", updated.Cover, updated.Location)
	}
	if !updated.IsTop || !updated.IsEssence || !updated.IsOutdated {
		t.Fatalf("updated flags = top %t essence %t outdated %t, want all true", updated.IsTop, updated.IsEssence, updated.IsOutdated)
	}
	if updated.CategoryID == nil || *updated.CategoryID != category.ID {
		t.Fatalf("category_id = %v, want %d", updated.CategoryID, category.ID)
	}
	assertMCPTagIDs(t, updated.Tags, []uint{tags[0].ID, tags[1].ID})

	draftUpdate, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name: "article_update_published",
		Arguments: map[string]any{
			"id":      draft.ID,
			"summary": "must not persist",
		},
	})
	if err != nil {
		t.Fatalf("call article_update_published for draft: %v", err)
	}
	if !draftUpdate.IsError {
		t.Fatal("article_update_published draft result isError = false, want true")
	}
	var draftAfter model.Article
	if err := db.First(&draftAfter, draft.ID).Error; err != nil {
		t.Fatalf("reload draft: %v", err)
	}
	if draftAfter.Summary != draft.Summary {
		t.Fatalf("draft summary changed to %q, want %q", draftAfter.Summary, draft.Summary)
	}

	unpublishResult := callMCPTool(t, ctx, session, "article_unpublish", map[string]any{"id": published.ID})
	var unpublishOutput tools.ArticleUnpublishOutput
	decodeStructuredContent(t, unpublishResult.StructuredContent, &unpublishOutput)
	if unpublishOutput.AlreadyUnpublished {
		t.Fatal("first article_unpublish returned already_unpublished=true")
	}
	if unpublishOutput.Item.IsPublish {
		t.Fatal("first article_unpublish returned is_publish=true")
	}

	unpublished := loadMCPArticle(t, db, published.ID)
	if unpublished.IsPublish {
		t.Fatal("article_unpublish persisted is_publish=true")
	}
	if unpublished.Title != updated.Title || unpublished.Content != updated.Content || unpublished.Summary != updated.Summary {
		t.Fatalf("unpublish changed content fields: title=%q content=%q summary=%q", unpublished.Title, unpublished.Content, unpublished.Summary)
	}
	if !unpublished.IsTop || !unpublished.IsEssence || !unpublished.IsOutdated {
		t.Fatalf("unpublish changed flags: top=%t essence=%t outdated=%t", unpublished.IsTop, unpublished.IsEssence, unpublished.IsOutdated)
	}

	secondUnpublish := callMCPTool(t, ctx, session, "article_unpublish", map[string]any{"id": published.ID})
	var secondUnpublishOutput tools.ArticleUnpublishOutput
	decodeStructuredContent(t, secondUnpublish.StructuredContent, &secondUnpublishOutput)
	if !secondUnpublishOutput.AlreadyUnpublished {
		t.Fatal("second article_unpublish returned already_unpublished=false, want true")
	}

	deleteResult := callMCPTool(t, ctx, session, "article_delete", map[string]any{"id": deletable.ID})
	var deleteOutput tools.ArticleDeleteOutput
	decodeStructuredContent(t, deleteResult.StructuredContent, &deleteOutput)
	if !deleteOutput.Deleted || deleteOutput.ID != deletable.ID {
		t.Fatalf("article_delete output = %+v", deleteOutput)
	}
	var deleted model.Article
	if err := db.First(&deleted, deletable.ID).Error; !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("deleted article lookup error = %v, want record not found", err)
	}

	secondDelete, err := session.CallTool(ctx, &sdkmcp.CallToolParams{
		Name:      "article_delete",
		Arguments: map[string]any{"id": deletable.ID},
	})
	if err != nil {
		t.Fatalf("call article_delete second time: %v", err)
	}
	if !secondDelete.IsError {
		t.Fatal("second article_delete result isError = false, want true")
	}
}
