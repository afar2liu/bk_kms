package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
)

var httpClient = &http.Client{
	Timeout: 60 * time.Second,
}

// BookmarkContent 书签内容
type BookmarkContent struct {
	Title   string
	Author  string
	Excerpt string
	Content string // 纯文本内容
	HTML    string // HTML 内容
}

// FetchBookmarkContent 获取书签内容
func FetchBookmarkContent(bookmarkURL string, keepTitle, keepExcerpt bool) (*BookmarkContent, error) {
	// 1. 下载网页内容
	req, err := http.NewRequest("GET", bookmarkURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	// 发送请求
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("下载网页失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 状态码错误: %d", resp.StatusCode)
	}

	// 获取 Content-Type
	contentType := resp.Header.Get("Content-Type")

	// 2. 解析网页内容（仅处理 HTML）
	result := &BookmarkContent{}

	if !strings.Contains(contentType, "text/html") {
		// 非 HTML 内容，只返回 URL 作为标题
		result.Title = bookmarkURL
		return result, nil
	}

	// 3. 使用 readability 解析文章
	parsedURL, err := url.Parse(bookmarkURL)
	if err != nil {
		return nil, fmt.Errorf("解析 URL 失败: %w", err)
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		return nil, fmt.Errorf("解析文章失败: %w", err)
	}

	// 4. 提取内容
	result.Author = article.Byline
	result.Content = article.TextContent
	result.HTML = article.Content

	// 如果不保留标题或标题为空，使用解析的标题
	if !keepTitle || article.Title != "" {
		result.Title = article.Title
	}

	// 如果不保留摘要或摘要为空，使用解析的摘要
	if !keepExcerpt || article.Excerpt != "" {
		result.Excerpt = article.Excerpt
	}

	// 确保标题不为空
	if result.Title == "" {
		result.Title = bookmarkURL
	}

	return result, nil
}

// RemoveUTMParams 移除 URL 中的 UTM 参数
func RemoveUTMParams(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// 获取查询参数
	queries := parsedURL.Query()

	// 移除 UTM 参数
	utmParams := []string{
		"utm_source", "utm_medium", "utm_campaign",
		"utm_term", "utm_content", "utm_name",
	}

	for _, param := range utmParams {
		queries.Del(param)
	}

	// 重新设置查询参数
	parsedURL.RawQuery = queries.Encode()

	return parsedURL.String(), nil
}

// ValidateTitle 验证标题是否为有效的 UTF-8，如果无效则使用 URL
func ValidateTitle(title, url string) string {
	if !utf8.ValidString(title) {
		return url
	}
	return title
}

// NormalizeSpace 标准化空白字符
func NormalizeSpace(s string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

// ParsedBookmark 解析后的书签
type ParsedBookmark struct {
	URL        string
	Title      string
	Tags       []string
	Category   string
	ModifiedAt time.Time
}

// ParseNetscapeBookmarkHTML 解析 Netscape Bookmark 格式的 HTML 文件
func ParseNetscapeBookmarkHTML(reader io.Reader, generateTag bool) ([]ParsedBookmark, error) {
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	var bookmarks []ParsedBookmark
	mapURL := make(map[string]struct{})

	doc.Find("dt>a").Each(func(_ int, a *goquery.Selection) {
		// 获取相关元素
		dt := a.Parent()
		dl := dt.Parent()
		h3 := dl.Parent().Find("h3").First()

		// 获取元数据
		title := a.Text()
		rawURL, _ := a.Attr("href")
		strTags, _ := a.Attr("tags")

		dateStr, fieldExists := a.Attr("last_modified")
		if !fieldExists {
			dateStr, _ = a.Attr("add_date")
		}

		// 使用当前时间作为默认日期
		modifiedDate := time.Now()
		if dateStr != "" {
			modifiedTsInt, err := strconv.ParseInt(dateStr, 10, 64)
			if err == nil {
				modifiedDate = time.Unix(modifiedTsInt, 0)
			}
		}

		// 清理 URL
		cleanURL, err := RemoveUTMParams(rawURL)
		if err != nil {
			// URL 无效，跳过
			return
		}

		// 验证标题
		title = ValidateTitle(title, cleanURL)

		// 检查 URL 是否已存在
		if _, exist := mapURL[cleanURL]; exist {
			return
		}

		// 获取标签
		var tags []string
		for _, strTag := range strings.Split(strTags, ",") {
			strTag = NormalizeSpace(strTag)
			if strTag != "" {
				tags = append(tags, strTag)
			}
		}

		// 获取分类名称并添加为标签（如果需要）
		category := NormalizeSpace(h3.Text())
		if category != "" && generateTag {
			tags = append(tags, category)
		}

		// 添加到列表
		bookmark := ParsedBookmark{
			URL:        cleanURL,
			Title:      title,
			Tags:       tags,
			Category:   category,
			ModifiedAt: modifiedDate,
		}

		mapURL[cleanURL] = struct{}{}
		bookmarks = append(bookmarks, bookmark)
	})

	return bookmarks, nil
}
