package tools

import "testing"

func TestNormalizeUserListPageIsBounded(t *testing.T) {
	tests := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "defaults", page: 0, pageSize: 0, wantPage: 1, wantPageSize: defaultUserListPageSize},
		{name: "negative values", page: -2, pageSize: -3, wantPage: 1, wantPageSize: defaultUserListPageSize},
		{name: "preserve valid", page: 3, pageSize: 50, wantPage: 3, wantPageSize: 50},
		{name: "cap oversized", page: 2, pageSize: 1000, wantPage: 2, wantPageSize: maxUserListPageSize},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			page, pageSize := normalizeUserListPage(tc.page, tc.pageSize)
			if page != tc.wantPage || pageSize != tc.wantPageSize {
				t.Fatalf("normalizeUserListPage(%d, %d) = (%d, %d), want (%d, %d)", tc.page, tc.pageSize, page, pageSize, tc.wantPage, tc.wantPageSize)
			}
		})
	}
}
