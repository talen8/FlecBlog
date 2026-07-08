package tools

import (
	"encoding/json"
	"testing"
)

func TestAggregateListOutputsPreserveEmptyListAndZeroTotal(t *testing.T) {
	commentList, commentTotal := ExplicitListFields([]CommentItem(nil), 0)
	taxonomyList, taxonomyTotal := ExplicitListFields([]TaxonomyItem(nil), 0)
	friendList, friendTotal := ExplicitListFields([]FriendItem(nil), 0)
	rssList, rssTotal := ExplicitListFields([]RssFeedItem(nil), 0)
	momentList, momentTotal := ExplicitListFields([]MomentItem(nil), 0)

	cases := []struct {
		name  string
		value any
	}{
		{"comment", CommentManageOutput{List: commentList, Total: commentTotal, Page: 1, PageSize: 10}},
		{"taxonomy", TaxonomyManageOutput{List: taxonomyList, Total: taxonomyTotal, Page: 1, PageSize: 10}},
		{"friend", FriendManageOutput{List: friendList, Total: friendTotal, Page: 1, PageSize: 10}},
		{"rssfeed", RssFeedManageOutput{List: rssList, Total: rssTotal, Page: 1, PageSize: 10}},
		{"moment", MomentManageOutput{List: momentList, Total: momentTotal, Page: 1, PageSize: 10}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			list, ok := got["list"].([]any)
			if !ok || len(list) != 0 {
				t.Fatalf("list = %#v, want []", got["list"])
			}
			if total, ok := got["total"].(float64); !ok || total != 0 {
				t.Fatalf("total = %#v, want 0", got["total"])
			}
		})
	}
}

func TestAggregateNonListOutputsOmitListFields(t *testing.T) {
	cases := []struct {
		name  string
		value any
	}{
		{"comment", CommentManageOutput{}},
		{"taxonomy", TaxonomyManageOutput{}},
		{"friend", FriendManageOutput{}},
		{"rssfeed", RssFeedManageOutput{}},
		{"moment", MomentManageOutput{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.value)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var got map[string]any
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if _, ok := got["list"]; ok {
				t.Fatalf("non-list output unexpectedly contains list: %s", data)
			}
			if _, ok := got["total"]; ok {
				t.Fatalf("non-list output unexpectedly contains total: %s", data)
			}
		})
	}
}
