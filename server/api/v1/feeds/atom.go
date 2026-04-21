package feeds

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"
	"time"

	"flec_blog/config"
	"flec_blog/internal/dto"
	"flec_blog/internal/service"

	"github.com/gin-gonic/gin"
)

// AtomController Atom控制器
type AtomController struct {
	articleService *service.ArticleService
	config         *config.Config
}

// NewAtomController 创建Atom控制器
func NewAtomController(articleService *service.ArticleService, config *config.Config) *AtomController {
	return &AtomController{
		articleService: articleService,
		config:         config,
	}
}

// GetAtomFeed 获取Atom订阅
//
//	@Summary		Atom订阅
//	@Description	生成博客文章的Atom 1.0订阅源
//	@Tags			订阅
//	@Accept			json
//	@Produce		xml
//	@Success		200	{string}	string	"Atom XML 订阅内容"
//	@Router			/atom.xml [get]
func (c *AtomController) GetAtomFeed(ctx *gin.Context) {
	req := &dto.ListArticlesForWebRequest{
		Page:     1,
		PageSize: 0,
	}

	articles, _, err := c.articleService.ListForWeb(ctx.Request.Context(), req)
	if err != nil {
		ctx.XML(500, gin.H{"error": "生成Atom失败"})
		return
	}

	baseURL := strings.TrimRight(c.config.Basic.BlogURL, "/")
	canonicalURL := ""
	if baseURL != "" {
		canonicalURL = baseURL + "/"
	}
	siteName := c.config.Blog.Title
	fallbackTime := time.Now().UTC().Truncate(time.Second)

	sort.SliceStable(articles, func(i, j int) bool {
		updatedI := fallbackTime
		if articles[i].UpdateTime != nil && !articles[i].UpdateTime.IsZero() {
			updatedI = articles[i].UpdateTime.Time.UTC().Truncate(time.Second)
		} else if articles[i].PublishTime != nil && !articles[i].PublishTime.IsZero() {
			updatedI = articles[i].PublishTime.Time.UTC().Truncate(time.Second)
		}

		updatedJ := fallbackTime
		if articles[j].UpdateTime != nil && !articles[j].UpdateTime.IsZero() {
			updatedJ = articles[j].UpdateTime.Time.UTC().Truncate(time.Second)
		} else if articles[j].PublishTime != nil && !articles[j].PublishTime.IsZero() {
			updatedJ = articles[j].PublishTime.Time.UTC().Truncate(time.Second)
		}

		return updatedI.After(updatedJ)
	})

	feedUpdatedAt := fallbackTime
	if len(articles) > 0 {
		if articles[0].UpdateTime != nil && !articles[0].UpdateTime.IsZero() {
			feedUpdatedAt = articles[0].UpdateTime.Time.UTC().Truncate(time.Second)
		} else if articles[0].PublishTime != nil && !articles[0].PublishTime.IsZero() {
			feedUpdatedAt = articles[0].PublishTime.Time.UTC().Truncate(time.Second)
		}
	}

	atom := &Atom{
		XMLNS:   "http://www.w3.org/2005/Atom",
		Title:   siteName,
		ID:      canonicalURL,
		Updated: feedUpdatedAt.Format(time.RFC3339),
		Author: &AtomAuthor{
			Name: siteName,
		},
		Link: []AtomLink{
			{Href: canonicalURL, Rel: "alternate"},
			{Href: fmt.Sprintf("%s/atom.xml", baseURL), Rel: "self"},
		},
		Entries: make([]AtomEntry, 0, len(articles)),
	}

	for _, article := range articles {
		articleURL := article.URL
		if baseURL != "" {
			articleURL = baseURL + article.URL
		}

		links := []AtomLink{{Href: articleURL}}
		if article.Cover != "" {
			links = append(links, AtomLink{
				Href: article.Cover,
				Rel:  "enclosure",
			})
		}

		entry := AtomEntry{
			Title:   article.Title,
			ID:      articleURL,
			Updated: feedUpdatedAt.Format(time.RFC3339),
			Link:    links,
			Summary: article.Summary,
		}

		if article.UpdateTime != nil && !article.UpdateTime.IsZero() {
			entry.Updated = article.UpdateTime.Time.UTC().Truncate(time.Second).Format(time.RFC3339)
		} else if article.PublishTime != nil && !article.PublishTime.IsZero() {
			entry.Updated = article.PublishTime.Time.UTC().Truncate(time.Second).Format(time.RFC3339)
		}

		if article.PublishTime != nil && !article.PublishTime.IsZero() {
			entry.Published = article.PublishTime.Time.UTC().Truncate(time.Second).Format(time.RFC3339)
		} else if article.UpdateTime != nil && !article.UpdateTime.IsZero() {
			entry.Published = article.UpdateTime.Time.UTC().Truncate(time.Second).Format(time.RFC3339)
		}

		if article.Category.Name != "" {
			entry.Category = []AtomCategory{{Term: article.Category.Name}}
		}

		atom.Entries = append(atom.Entries, entry)
	}

	xmlData, err := xml.MarshalIndent(atom, "", "  ")
	if err != nil {
		ctx.XML(500, gin.H{"error": "生成Atom失败"})
		return
	}

	ctx.Header("Content-Type", "application/atom+xml; charset=utf-8")
	ctx.String(200, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+string(xmlData))
}

// Atom Feed 结构定义
type Atom struct {
	XMLName xml.Name    `xml:"feed"`
	XMLNS   string      `xml:"xmlns,attr"`
	Title   string      `xml:"title"`
	ID      string      `xml:"id"`
	Updated string      `xml:"updated"`
	Author  *AtomAuthor `xml:"author,omitempty"`
	Link    []AtomLink  `xml:"link"`
	Entries []AtomEntry `xml:"entry"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
}

type AtomEntry struct {
	Title     string         `xml:"title"`
	ID        string         `xml:"id"`
	Updated   string         `xml:"updated"`
	Published string         `xml:"published,omitempty"`
	Link      []AtomLink     `xml:"link"`
	Summary   string         `xml:"summary"`
	Category  []AtomCategory `xml:"category,omitempty"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
}

type AtomCategory struct {
	Term string `xml:"term,attr"`
}
