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

// RSSController RSS 2.0控制器
type RSSController struct {
	articleService *service.ArticleService
	config         *config.Config
}

// NewRSSController 创建RSS 2.0控制器
func NewRSSController(articleService *service.ArticleService, config *config.Config) *RSSController {
	return &RSSController{
		articleService: articleService,
		config:         config,
	}
}

// GetRSSFeed 获取RSS 2.0订阅
//
//	@Summary		RSS 2.0订阅
//	@Description	生成博客文章的RSS 2.0订阅源
//	@Tags			订阅
//	@Accept			json
//	@Produce		xml
//	@Success		200	{string}	string	"RSS XML 订阅内容"
//	@Router			/rss.xml [get]
func (c *RSSController) GetRSSFeed(ctx *gin.Context) {
	req := &dto.ListArticlesForWebRequest{
		Page:     1,
		PageSize: 0,
	}

	articles, _, err := c.articleService.ListForWeb(ctx.Request.Context(), req)
	if err != nil {
		ctx.XML(500, gin.H{"error": "生成RSS失败"})
		return
	}

	baseURL := strings.TrimRight(c.config.Basic.BlogURL, "/")
	canonicalURL := ""
	if baseURL != "" {
		canonicalURL = baseURL + "/"
	}
	siteName := c.config.Blog.Title
	siteDescription := c.config.Blog.Description
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

	rss := &RSS{
		Version: "2.0",
		AtomNS:  "http://www.w3.org/2005/Atom",
		Channel: RSSChannel{
			Title:         siteName,
			Link:          canonicalURL,
			Description:   siteDescription,
			Language:      "zh-CN",
			LastBuildDate: feedUpdatedAt.Format(time.RFC1123Z),
			AtomLink: RSSAtomLink{
				Href: fmt.Sprintf("%s/rss.xml", baseURL),
				Rel:  "self",
				Type: "application/rss+xml",
			},
			Items: make([]RSSItem, 0, len(articles)),
		},
	}

	for _, article := range articles {
		articleURL := article.URL
		if baseURL != "" {
			articleURL = baseURL + article.URL
		}
		item := RSSItem{
			Title:       article.Title,
			Link:        articleURL,
			GUID:        RSSGUID{IsPermaLink: "true", Value: articleURL},
			Description: article.Summary,
		}

		if article.PublishTime != nil && !article.PublishTime.IsZero() {
			item.PubDate = article.PublishTime.Time.UTC().Truncate(time.Second).Format(time.RFC1123Z)
		} else if article.UpdateTime != nil && !article.UpdateTime.IsZero() {
			item.PubDate = article.UpdateTime.Time.UTC().Truncate(time.Second).Format(time.RFC1123Z)
		}

		if article.Category.Name != "" {
			item.Category = article.Category.Name
		}

		rss.Channel.Items = append(rss.Channel.Items, item)
	}

	xmlData, err := xml.MarshalIndent(rss, "", "  ")
	if err != nil {
		ctx.XML(500, gin.H{"error": "生成RSS失败"})
		return
	}

	ctx.Header("Content-Type", "application/rss+xml; charset=utf-8")
	ctx.String(200, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"+string(xmlData))
}

// RSS 2.0 结构定义
type RSS struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	AtomNS  string     `xml:"xmlns:atom,attr,omitempty"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title         string      `xml:"title"`
	Link          string      `xml:"link"`
	Description   string      `xml:"description"`
	Language      string      `xml:"language,omitempty"`
	LastBuildDate string      `xml:"lastBuildDate,omitempty"`
	AtomLink      RSSAtomLink `xml:"atom:link"`
	Items         []RSSItem   `xml:"item"`
}

type RSSAtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr"`
	Type string `xml:"type,attr"`
}

type RSSItem struct {
	Title       string  `xml:"title"`
	Link        string  `xml:"link"`
	GUID        RSSGUID `xml:"guid"`
	Description string  `xml:"description"`
	PubDate     string  `xml:"pubDate,omitempty"`
	Category    string  `xml:"category,omitempty"`
}

type RSSGUID struct {
	IsPermaLink string `xml:"isPermaLink,attr,omitempty"`
	Value       string `xml:",chardata"`
}
