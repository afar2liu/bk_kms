package repo

import (
	"bk_kms/lib"
	"bk_kms/model/db"
	"strings"

	"gorm.io/gorm"
)

type BookmarkRepo struct{}

// List 查询书签列表
func (r *BookmarkRepo) List(keyword string, tags []string, page, pageSize int) ([]db.Bookmark, int64, error) {
	var bookmarks []db.Bookmark
	var total int64

	query := lib.DB.Model(&db.Bookmark{})

	// 关键字搜索
	if keyword != "" {
		query = query.Where("url LIKE ? OR title LIKE ? OR excerpt LIKE ? OR content LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 标签过滤
	if len(tags) > 0 {
		query = query.Joins("JOIN bookmark_tag ON bookmark.id = bookmark_tag.bookmark_id").
			Joins("JOIN tag ON bookmark_tag.tag_id = tag.id").
			Where("tag.name IN ?", tags).
			Group("bookmark.id").
			Having("COUNT(DISTINCT tag.id) = ?", len(tags))
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Tags").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bookmarks).Error

	return bookmarks, total, err
}

// FindByID 根据ID查找书签
func (r *BookmarkRepo) FindByID(id int) (*db.Bookmark, error) {
	var bookmark db.Bookmark
	err := lib.DB.Preload("Tags").First(&bookmark, id).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

// FindByURL 根据URL查找书签
func (r *BookmarkRepo) FindByURL(url string) (*db.Bookmark, error) {
	var bookmark db.Bookmark
	err := lib.DB.Where("url = ?", url).First(&bookmark).Error
	if err != nil {
		return nil, err
	}
	return &bookmark, nil
}

// Create 创建书签
func (r *BookmarkRepo) Create(bookmark *db.Bookmark) error {
	return lib.DB.Create(bookmark).Error
}

// Update 更新书签
func (r *BookmarkRepo) Update(bookmark *db.Bookmark) error {
	return lib.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(bookmark).Error
}

// Delete 删除书签
func (r *BookmarkRepo) Delete(ids []int) error {
	return lib.DB.Transaction(func(tx *gorm.DB) error {
		// 删除关联的标签
		if err := tx.Where("bookmark_id IN ?", ids).Delete(&db.BookmarkTag{}).Error; err != nil {
			return err
		}
		// 删除书签
		return tx.Delete(&db.Bookmark{}, ids).Error
	})
}

// FindOrCreateTags 查找或创建标签
func (r *BookmarkRepo) FindOrCreateTags(tagNames []string) ([]db.Tag, error) {
	var tags []db.Tag

	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		var tag db.Tag
		err := lib.DB.Where("name = ?", name).First(&tag).Error
		if err == gorm.ErrRecordNotFound {
			// 创建新标签
			tag = db.Tag{Name: name}
			if err := lib.DB.Create(&tag).Error; err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	return tags, nil
}
