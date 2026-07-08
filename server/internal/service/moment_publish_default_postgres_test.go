package service

import (
	"context"
	"testing"

	"flec_blog/internal/dto"
	"flec_blog/internal/model"
	"flec_blog/internal/repository"
	"flec_blog/internal/testutil"
)

func TestMomentCreatePublishTriStateOverPostgres(t *testing.T) {
	db := testutil.OpenMCPTestPostgres(t)

	service := NewMomentService(repository.NewMomentRepository(db), nil)
	ctx := context.Background()

	falseValue := false
	trueValue := true
	tests := []struct {
		name      string
		isPublish *bool
		want      bool
	}{
		{name: "omitted defaults true", isPublish: nil, want: true},
		{name: "explicit false preserved", isPublish: &falseValue, want: false},
		{name: "explicit true preserved", isPublish: &trueValue, want: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			created, err := service.Create(ctx, &dto.CreateMomentRequest{
				Content:   dto.MomentContent{Text: "moment publish tri-state regression"},
				IsPublish: tc.isPublish,
			})
			if err != nil {
				t.Fatalf("create moment: %v", err)
			}
			t.Cleanup(func() {
				_ = db.Delete(&model.Moment{}, created.ID).Error
			})
			if created.IsPublish != tc.want {
				t.Fatalf("created is_publish = %v, want %v", created.IsPublish, tc.want)
			}

			var stored model.Moment
			if err := db.First(&stored, created.ID).Error; err != nil {
				t.Fatalf("reload moment: %v", err)
			}
			if stored.IsPublish != tc.want {
				t.Fatalf("stored is_publish = %v, want %v", stored.IsPublish, tc.want)
			}
		})
	}
}
