package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"bk_kms/lib"
	"bk_kms/model/dto"
	"bk_kms/repo"
)

type TagController struct {
	tagRepo *repo.TagRepo
}

func NewTagController() *TagController {
	return &TagController{
		tagRepo: &repo.TagRepo{},
	}
}

// List tag列表
func (tc *TagController) List(c *gin.Context) {
	var req dto.TagListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	tags, err := tc.tagRepo.List(req.Name)
	if err != nil {
		lib.Logger.Error("查询标签列表失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "查询失败",
		})
		return
	}

	// 转换为 DTO
	items := make([]dto.TagListItem, 0, len(tags))
	for _, tag := range tags {
		items = append(items, dto.TagListItem{
			ID:    tag.ID,
			Name:  tag.Name,
			Count: tag.Count,
		})
	}

	c.JSON(http.StatusOK, dto.TagListResponse{
		Code: 0,
		Msg:  "成功",
		Data: items,
	})
}

// Update tag重命名
func (tc *TagController) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误",
		})
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "参数错误: " + err.Error(),
		})
		return
	}

	// 查找标签
	tag, err := tc.tagRepo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, dto.Response{
				Code: 1,
				Msg:  "标签不存在",
			})
			return
		}
		lib.Logger.Error("查询标签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "更新失败",
		})
		return
	}

	// 更新标签
	tag.Name = req.Name
	if err := tc.tagRepo.Update(tag); err != nil {
		lib.Logger.Error("更新标签失败: " + err.Error())
		c.JSON(http.StatusOK, dto.Response{
			Code: 1,
			Msg:  "更新失败",
		})
		return
	}

	lib.Logger.Info("更新标签成功: " + req.Name)

	c.JSON(http.StatusOK, dto.Response{
		Code: 0,
		Msg:  "更新成功",
		Data: struct{}{},
	})
}
