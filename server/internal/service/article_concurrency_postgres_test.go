package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/testutil"
)

func TestArticleOptimisticLockOverPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	articleRepo := repository.NewArticleRepository(db)
	articleService := NewArticleService(
		articleRepo,
		repository.NewTagRepository(db),
		repository.NewCategoryRepository(db),
		repository.NewCommentRepository(db),
		nil,
		db,
	)
	ctx := context.Background()

	article := &model.Article{
		Title:     "optimistic-lock",
		Slug:      fmt.Sprintf("optimistic-lock-%d", time.Now().UnixNano()),
		Content:   "original content",
		Summary:   "original summary",
		IsPublish: false,
	}
	if err := db.Create(article).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}
	t.Cleanup(func() { _ = db.Unscoped().Delete(&model.Article{}, article.ID).Error })

	first, err := articleService.Get(ctx, article.ID)
	if err != nil {
		t.Fatalf("first snapshot: %v", err)
	}
	second, err := articleService.Get(ctx, article.ID)
	if err != nil {
		t.Fatalf("second snapshot: %v", err)
	}
	if first.Revision.IsZero() || !first.Revision.Equal(second.Revision) {
		t.Fatalf("snapshot revisions first=%s second=%s", first.Revision, second.Revision)
	}

	firstReq := updateRequestFromSnapshot(first)
	firstReq.Summary = "winner"
	firstReq.ExpectedUpdatedAt = timePtr(first.Revision)
	updated, err := articleService.Update(ctx, article.ID, firstReq)
	if err != nil {
		t.Fatalf("first update: %v", err)
	}
	if !updated.Revision.After(first.Revision) {
		t.Fatalf("revision did not advance: before=%s after=%s", first.Revision, updated.Revision)
	}

	secondReq := updateRequestFromSnapshot(second)
	secondReq.Summary = "stale overwrite"
	secondReq.ExpectedUpdatedAt = timePtr(second.Revision)
	if _, err := articleService.Update(ctx, article.ID, secondReq); !errors.Is(err, ErrArticleUpdateConflict) {
		t.Fatalf("stale update error = %v, want ErrArticleUpdateConflict", err)
	}

	final, err := articleService.Get(ctx, article.ID)
	if err != nil {
		t.Fatalf("final article: %v", err)
	}
	if final.Summary != "winner" {
		t.Fatalf("final summary = %q, want winner", final.Summary)
	}
}

func TestConcurrentArticlePublishOnlyOneUpdateWinsOverPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	articleService := NewArticleService(
		repository.NewArticleRepository(db),
		repository.NewTagRepository(db),
		repository.NewCategoryRepository(db),
		repository.NewCommentRepository(db),
		nil,
		db,
	)
	ctx := context.Background()
	article := &model.Article{
		Title:     "concurrent-publish",
		Slug:      fmt.Sprintf("concurrent-publish-%d", time.Now().UnixNano()),
		Content:   "content",
		Summary:   "summary",
		IsPublish: false,
	}
	if err := db.Create(article).Error; err != nil {
		t.Fatalf("create article: %v", err)
	}
	t.Cleanup(func() { _ = db.Unscoped().Delete(&model.Article{}, article.ID).Error })

	left, err := articleService.Get(ctx, article.ID)
	if err != nil {
		t.Fatalf("left snapshot: %v", err)
	}
	right, err := articleService.Get(ctx, article.ID)
	if err != nil {
		t.Fatalf("right snapshot: %v", err)
	}
	publish := true
	leftReq := updateRequestFromSnapshot(left)
	leftReq.IsPublish = &publish
	leftReq.ExpectedUpdatedAt = timePtr(left.Revision)
	rightReq := updateRequestFromSnapshot(right)
	rightReq.IsPublish = &publish
	rightReq.ExpectedUpdatedAt = timePtr(right.Revision)

	start := make(chan struct{})
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, req := range []*dto.UpdateArticleRequest{leftReq, rightReq} {
		req := req
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, err := articleService.Update(ctx, article.ID, req)
			results <- err
		}()
	}
	close(start)
	wg.Wait()
	close(results)

	successes := 0
	conflicts := 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrArticleUpdateConflict):
			conflicts++
		default:
			t.Fatalf("unexpected publish error: %v", err)
		}
	}
	if successes != 1 || conflicts != 1 {
		t.Fatalf("publish results successes=%d conflicts=%d", successes, conflicts)
	}
	final, err := articleService.Get(ctx, article.ID)
	if err != nil {
		t.Fatalf("final article: %v", err)
	}
	if !final.IsPublish {
		t.Fatal("final article is not published")
	}
}

func updateRequestFromSnapshot(current *dto.ArticleAdminDetailResponse) *dto.UpdateArticleRequest {
	return &dto.UpdateArticleRequest{
		Summary:    current.Summary,
		AISummary:  current.AISummary,
		Cover:      current.Cover,
		Location:   current.Location,
		CategoryID: articleSnapshotCategoryID(current.Category.ID),
	}
}

func articleSnapshotCategoryID(id uint) *uint {
	if id == 0 {
		return nil
	}
	value := id
	return &value
}

func timePtr(value time.Time) *time.Time {
	copy := value
	return &copy
}
