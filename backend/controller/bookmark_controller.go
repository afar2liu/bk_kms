package controller

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bk_kms/lib"
	"bk_kms/model/db"
	"bk_kms/model/dto"
	"bk_kms/repo"
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
	// TODO: 根据URL获取书签内容


	// 创建书签
	bookmark := &db.Bookmark{
		URL:       req.URL,
		Title:     req.Title,
		Excerpt:   req.Excerpt,
		Author:    "",
		Content:   "",
		HTML:      "",
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
