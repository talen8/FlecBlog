package tools

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/service"
	"flec_blog/internal/testutil"

	"gorm.io/gorm"
)

func TestConcurrentPublishToolIsIdempotentOverPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	articleService := service.NewArticleService(
		repository.NewArticleRepository(db),
		repository.NewTagRepository(db),
		repository.NewCategoryRepository(db),
		repository.NewCommentRepository(db),
		nil,
		db,
	)
	wrapper := NewArticleWrapper(articleService)
	article := &model.Article{
		Title:     "tool-concurrent-publish",
		Slug:      fmt.Sprintf("tool-concurrent-publish-%d", time.Now().UnixNano()),
		Content:   "content",
		Summary:   "summary",
		IsPublish: false,
	}
	if err := db.Create(article).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}
	t.Cleanup(func() { _ = db.Unscoped().Delete(&model.Article{}, article.ID).Error })

	// ---------- test-only update barrier ----------
	// Forces both concurrent Publish updates to reach the DB layer before either
	// can commit, guaranteeing they race on the same captured pre-publish revision.
	const callbackName = "test:article_publish_conflict_barrier"
	arrived := make(chan struct{}, 2)
	release := make(chan struct{})

	var releaseOnce sync.Once
	releaseBarrier := func() { releaseOnce.Do(func() { close(release) }) }
	t.Cleanup(releaseBarrier)

	if err := db.Callback().Update().
		Before("gorm:update").
		Register(callbackName, func(tx *gorm.DB) {
			if tx.Statement.Table != "articles" {
				return
			}
			// Narrow to target article when Dest is available.
			if a, ok := tx.Statement.Dest.(*model.Article); ok && a.ID != article.ID {
				return
			}
			arrived <- struct{}{}
			<-release
		}); err != nil {
		t.Fatalf("register update barrier: %v", err)
	}
	t.Cleanup(func() { _ = db.Callback().Update().Remove(callbackName) })
	// ---------- end barrier ----------

	type outcome struct {
		output ArticlePublishOutput
		err    error
	}
	start := make(chan struct{})
	outcomes := make(chan outcome, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, output, err := wrapper.PublishArticleTool(context.Background(), nil, ArticlePublishInput{ID: article.ID})
			outcomes <- outcome{output: output, err: err}
		}()
	}
	close(start)

	// Wait for both goroutines to hit the update barrier.
	for i := 0; i < 2; i++ {
		select {
		case <-arrived:
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for concurrent article updates to reach barrier")
		}
	}
	releaseBarrier()
	wg.Wait()
	close(outcomes)

	successes := 0
	alreadyPublished := 0
	for result := range outcomes {
		if result.err != nil {
			t.Fatalf("publish tool error: %v", result.err)
		}
		successes++
		if result.output.AlreadyPublished {
			alreadyPublished++
		}
	}
	if successes != 2 || alreadyPublished != 1 {
		t.Fatalf("publish outcomes successes=%d already_published=%d", successes, alreadyPublished)
	}

	final, err := articleService.Get(context.Background(), article.ID)
	if err != nil {
		t.Fatalf("final article: %v", err)
	}
	if !final.IsPublish {
		t.Fatal("final article is not published")
	}
}
