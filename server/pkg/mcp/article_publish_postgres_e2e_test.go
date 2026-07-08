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

func TestArticlePublishOverMCPPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	category := model.Category{Name: "MCP Publish E2E Category", Slug: "mcp-publish-e2e-category"}
	if err := db.Create(&category).Error; err != nil {
		t.Fatalf("create category: %v", err)
	}
	tags := []model.Tag{
		{Name: "MCP Publish E2E Tag A", Slug: "mcp-publish-e2e-tag-a"},
		{Name: "MCP Publish E2E Tag B", Slug: "mcp-publish-e2e-tag-b"},
	}
	if err := db.Create(&tags).Error; err != nil {
		t.Fatalf("create tags: %v", err)
	}

	article := model.Article{
		Title:      "MCP Publish E2E Draft",
		Slug:       "mcp-publish-e2e-draft",
		Content:    "publish e2e body",
		Summary:    "publish e2e summary",
		AISummary:  "publish e2e ai summary",
		Cover:      "publish-cover.png",
		Location:   "Shenzhen",
		IsPublish:  false,
		CategoryID: &category.ID,
	}
	if err := db.Create(&article).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}
	if err := db.Model(&article).Association("Tags").Replace(tags); err != nil {
		t.Fatalf("attach tags: %v", err)
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

	publishResult := callMCPTool(t, ctx, session, "article_publish", map[string]any{"id": article.ID})
	var publishOutput tools.ArticlePublishOutput
	decodeStructuredContent(t, publishResult.StructuredContent, &publishOutput)
	if publishOutput.AlreadyPublished {
		t.Fatal("first article_publish returned already_published=true")
	}
	if !publishOutput.Item.IsPublish {
		t.Fatal("first article_publish returned is_publish=false")
	}

	published := loadMCPArticle(t, db, article.ID)
	if !published.IsPublish {
		t.Fatal("published article persisted is_publish=false")
	}
	if published.PublishTime == nil {
		t.Fatal("published article publish_time is nil")
	}
	if published.Title != article.Title || published.Content != article.Content {
		t.Fatalf("title/content changed unexpectedly: title=%q content=%q", published.Title, published.Content)
	}
	if published.Summary != article.Summary || published.AISummary != article.AISummary {
		t.Fatalf("summary fields changed unexpectedly: summary=%q ai_summary=%q", published.Summary, published.AISummary)
	}
	if published.Cover != article.Cover || published.Location != article.Location {
		t.Fatalf("metadata changed unexpectedly: cover=%q location=%q", published.Cover, published.Location)
	}
	if published.CategoryID == nil || *published.CategoryID != category.ID {
		t.Fatalf("category_id = %v, want %d", published.CategoryID, category.ID)
	}
	assertMCPTagIDs(t, published.Tags, []uint{tags[0].ID, tags[1].ID})
	firstPublishTime := *published.PublishTime

	secondResult := callMCPTool(t, ctx, session, "article_publish", map[string]any{"id": article.ID})
	var secondOutput tools.ArticlePublishOutput
	decodeStructuredContent(t, secondResult.StructuredContent, &secondOutput)
	if !secondOutput.AlreadyPublished {
		t.Fatal("second article_publish returned already_published=false, want true")
	}

	publishedAgain := loadMCPArticle(t, db, article.ID)
	if publishedAgain.PublishTime == nil || !publishedAgain.PublishTime.Equal(firstPublishTime) {
		t.Fatalf("repeat publish changed publish_time: got %v want %v", publishedAgain.PublishTime, firstPublishTime)
	}
}
