package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bk_kms/lib"
	"bk_kms/model/db"
	"bk_kms/model/dto"
	"bk_kms/repo"
	"bk_kms/utils"
)

type BookmarkController struct {
	bookmarkRepo *repo.BookmarkRepo
}

func NewBookmarkController() *BookmarkController {
	return &BookmarkController{
		bookmarkRepo: &repo.BookmarkRepo{},
	}
}

// List 书签列表
func (bc *BookmarkController) List(c *gin.Context) {
	var req dto.BookmarkListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认分页大小
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 解析 tags
	var tags []string
	if req.Tags != "" {
		tags = strings.Split(req.Tags, ",")
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	// 查询书签列表
	bookmarks, total, err := bc.bookmarkRepo.List(req.Keyword, tags, req.Page, req.PageSize)
	if err != nil {
		lib.Logger.Error("查询书签列表失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "查询失败",
		})
		return
	}

	// 转换为 DTO
	items := make([]dto.BookmarkListItem, 0, len(bookmarks))
	for _, bm := range bookmarks {
		tagItems := make([]dto.TagItem, 0, len(bm.Tags))
		for _, tag := range bm.Tags {
			tagItems = append(tagItems, dto.TagItem{
				ID:   tag.ID,
				Name: tag.Name,
			})
		}

		items = append(items, dto.BookmarkListItem{
			ID:        bm.ID,
			URL:       bm.URL,
			Title:     bm.Title,
			Excerpt:   bm.Excerpt,
			Author:    bm.Author,
			IsArchive: bm.IsArchive,
			CreatedAt: bm.CreatedAt.Unix(),
			UpdatedAt: bm.UpdatedAt.Unix(),
			Tags:      tagItems,
		})
	}

	c.JSON(http.StatusOK, dto.BookmarkListResponse{
		Code: 0,
		Msg:  "成功",
		Data: dto.PageData{
			Rows:  items,
			Total: int(total),
		},
	})
}

// Create 创建书签
func (bc *BookmarkController) Create(c *gin.Context) {
	var req dto.CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// 检查 URL 是否已存在
	if _, err := bc.bookmarkRepo.FindByURL(req.URL); err == nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "书签已存在",
		})
		return
	}

	// 查找或创建标签
	tagNames := make([]string, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	tags, err := bc.bookmarkRepo.FindOrCreateTags(tagNames)
	if err != nil {
		lib.Logger.Error("处理标签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "创建失败",
		})
		return
	}

	// 根据 URL 获取书签内容
	userHasDefinedTitle := req.Title != ""
	userHasDefinedExcerpt := req.Excerpt != ""

	var title, excerpt, author, content, html string
	title = req.Title
	excerpt = req.Excerpt

	// 如果需要创建归档，则获取网页内容
	if req.CreateArchive {
		lib.Logger.Info("开始获取书签内容: " + req.URL)
		bookmarkContent, err := utils.FetchBookmarkContent(req.URL, userHasDefinedTitle, userHasDefinedExcerpt)
		if err != nil {
			lib.Logger.Error("获取书签内容失败: " + err.Error())
			// 获取内容失败不影响创建书签，继续使用用户提供的信息
		} else {
			// 使用获取到的内容
			if !userHasDefinedTitle || title == "" {
				title = bookmarkContent.Title
			}
			if !userHasDefinedExcerpt || excerpt == "" {
				excerpt = bookmarkContent.Excerpt
			}
			author = bookmarkContent.Author
			content = bookmarkContent.Content
			html = bookmarkContent.HTML
			lib.Logger.Info("书签内容获取成功")
		}
	}

	// 确保标题不为空
	if title == "" {
		title = req.URL
	}

	// 创建书签
	bookmark := &db.Bookmark{
		URL:       req.URL,
		Title:     title,
		Excerpt:   excerpt,
		Author:    author,
		Content:   content,
		HTML:      html,
		IsArchive: req.CreateArchive,
		Tags:      tags,
	}

	if err := bc.bookmarkRepo.Create(bookmark); err != nil {
		lib.Logger.Error("创建书签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "创建失败",
		})
		return
	}

	lib.Logger.Info("创建书签成功: " + req.Title)

	c.JSON(http.StatusOK, dto.Response{
		Code: 0,
		Msg:  "创建成功",
	})
}

// Update 编辑书签
func (bc *BookmarkController) Update(c *gin.Context) {
	var req dto.UpdateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// 查找书签
	bookmark, err := bc.bookmarkRepo.FindByID(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, dto.Response{
				Code: 1,
				Msg:  "书签不存在",
			})
			return
		}
		lib.Logger.Error("查询书签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "更新失败",
		})
		return
	}

	// 查找或创建标签
	tagNames := make([]string, 0, len(req.Tags))
	for _, tag := range req.Tags {
		tagNames = append(tagNames, tag.Name)
	}
	tags, err := bc.bookmarkRepo.FindOrCreateTags(tagNames)
	if err != nil {
		lib.Logger.Error("处理标签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "更新失败",
		})
		return
	}

	// 更新书签
	bookmark.URL = req.URL
	bookmark.Title = req.Title
	bookmark.Excerpt = req.Excerpt
	bookmark.Author = req.Author
	bookmark.IsArchive = req.CreateArchive
	bookmark.Tags = tags
	// 如果需要创建归档，则获取网页内容
	if req.CreateArchive {
		lib.Logger.Info("开始获取书签内容: " + req.URL)
		bookmarkContent, err := utils.FetchBookmarkContent(req.URL, false, false)
		if err != nil {
			lib.Logger.Error("获取书签内容失败: " + err.Error())
			// 获取内容失败不影响创建书签，继续使用用户提供的信息
		} else {
			// 使用获取到的内容
			bookmark.Title = bookmarkContent.Title
			bookmark.Excerpt = bookmarkContent.Excerpt
			bookmark.Author = bookmarkContent.Author
			bookmark.Content = bookmarkContent.Content
			bookmark.HTML = bookmarkContent.HTML
			lib.Logger.Info("书签内容获取成功")
		}
	}

	if err := bc.bookmarkRepo.Update(bookmark); err != nil {
		lib.Logger.Error("更新书签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "更新失败",
		})
		return
	}

	lib.Logger.Info("更新书签成功: " + req.Title)

	c.JSON(http.StatusOK, dto.Response{
		Code: 0,
		Msg:  "更新成功",
	})
}

// Delete 删除书签
func (bc *BookmarkController) Delete(c *gin.Context) {
	var ids dto.DeleteBookmarkRequest
	if err := c.ShouldBindJSON(&ids); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	if len(ids) == 0 {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "请选择要删除的书签",
		})
		return
	}

	if err := bc.bookmarkRepo.Delete(ids); err != nil {
		lib.Logger.Error("删除书签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "删除失败",
		})
		return
	}

	lib.Logger.Info("删除书签成功")

	c.JSON(http.StatusOK, dto.Response{
		Code: 0,
		Msg:  "删除成功",
	})
}

// GetContent 查看内容
func (bc *BookmarkController) GetContent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误",
		})
		return
	}

	bookmark, err := bc.bookmarkRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, dto.Response{
				Code: 1,
				Msg:  "书签不存在",
			})
			return
		}
		lib.Logger.Error("查询书签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "查询失败",
		})
		return
	}

	c.JSON(http.StatusOK, dto.BookmarkContentResponse{
		Code: 0,
		Msg:  "成功",
		Data: dto.BookmarkContentData{
			ID:        bookmark.ID,
			URL:       bookmark.URL,
			Title:     bookmark.Title,
			HTML:      bookmark.HTML,
			CreatedAt: bookmark.CreatedAt.Unix(),
			UpdatedAt: bookmark.UpdatedAt.Unix(),
		},
	})
}

// Import 批量导入书签（使用 SSE 实时响应）
func (bc *BookmarkController) Import(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("bookmark_file")
	if err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "文件上传失败: " + err.Error(),
		})
		return
	}

	// 检查文件扩展名
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".html") {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "仅支持 .html 文件",
		})
		return
	}

	// 打开文件
	fileReader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "打开文件失败: " + err.Error(),
		})
		return
	}
	defer fileReader.Close()

	// 解析 HTML 文件
	generateTag := c.DefaultQuery("generate_tag", "false") == "true"
	bookmarks, err := utils.ParseNetscapeBookmarkHTML(fileReader, generateTag)
	if err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "解析书签文件失败: " + err.Error(),
		})
		return
	}

	if len(bookmarks) == 0 {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "未找到有效的书签",
		})
		return
	}

	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	// 发送 SSE 事件的辅助函数
	sendEvent := func(event dto.ImportProgressEvent) {
		data := fmt.Sprintf("data: %s\n\n", toJSON(event))
		c.Writer.WriteString(data)
		c.Writer.Flush()
	}

	// 发送开始事件
	sendEvent(dto.ImportProgressEvent{
		Type:    "progress",
		Message: fmt.Sprintf("开始导入，共 %d 个书签", len(bookmarks)),
		Current: 0,
		Total:   len(bookmarks),
	})

	// 统计信息
	successCount := 0
	skipCount := 0
	errorCount := 0

	// 逐个处理书签
	for i, bm := range bookmarks {
		// 检查 URL 是否已存在
		if _, err := bc.bookmarkRepo.FindByURL(bm.URL); err == nil {
			skipCount++
			sendEvent(dto.ImportProgressEvent{
				Type:    "progress",
				Message: fmt.Sprintf("跳过（已存在）: %s", bm.URL),
				Current: i + 1,
				Total:   len(bookmarks),
				URL:     bm.URL,
			})
			continue
		}

		// 查找或创建标签
		tags, err := bc.bookmarkRepo.FindOrCreateTags(bm.Tags)
		if err != nil {
			errorCount++
			sendEvent(dto.ImportProgressEvent{
				Type:    "error",
				Message: fmt.Sprintf("处理标签失败: %s - %v", bm.URL, err),
				Current: i + 1,
				Total:   len(bookmarks),
				URL:     bm.URL,
			})
			continue
		}

		// 创建书签
		bookmark := &db.Bookmark{
			URL:       bm.URL,
			Title:     bm.Title,
			Tags:      tags,
			IsArchive: false,
		}

		if err := bc.bookmarkRepo.Create(bookmark); err != nil {
			errorCount++
			sendEvent(dto.ImportProgressEvent{
				Type:    "error",
				Message: fmt.Sprintf("创建书签失败: %s - %v", bm.URL, err),
				Current: i + 1,
				Total:   len(bookmarks),
				URL:     bm.URL,
			})
			continue
		}

		successCount++
		sendEvent(dto.ImportProgressEvent{
			Type:    "success",
			Message: fmt.Sprintf("成功导入: %s", bm.Title),
			Current: i + 1,
			Total:   len(bookmarks),
			URL:     bm.URL,
		})
	}

	// 发送完成事件
	sendEvent(dto.ImportProgressEvent{
		Type:    "complete",
		Message: fmt.Sprintf("导入完成！成功: %d, 跳过: %d, 失败: %d", successCount, skipCount, errorCount),
		Current: len(bookmarks),
		Total:   len(bookmarks),
	})

	lib.Logger.Info(fmt.Sprintf("书签导入完成: 成功=%d, 跳过=%d, 失败=%d", successCount, skipCount, errorCount))
}

// toJSON 将对象转换为 JSON 字符串
func toJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}
