package dto

// PatchArticleRequest 文章局部更新请求
// nil 表示保留原值；非 nil 表示显式更新。
type PatchArticleRequest struct {
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	Summary    *string `json:"summary"`
	AISummary  *string `json:"ai_summary"`
	Cover      *string `json:"cover"`
	Location   *string `json:"location"`
	IsPublish  *bool   `json:"is_publish"`
	IsTop      *bool   `json:"is_top"`
	IsEssence  *bool   `json:"is_essence"`
	IsOutdated *bool   `json:"is_outdated"`
	CategoryID *uint   `json:"category_id"`
	TagIDs     []uint  `json:"tag_ids"`
}
