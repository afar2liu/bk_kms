package repo

import (
	"bk_kms/lib"
	"bk_kms/model/db"
)

type TagRepo struct{}

// TagWithCount Tag 带数量统计
type TagWithCount struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// List 查询标签列表
func (r *TagRepo) List(name string) ([]TagWithCount, error) {
	var tags []TagWithCount

	query := lib.DB.Table("tag").
		Select("tag.id, tag.name, COUNT(bookmark_tag.bookmark_id) as count").
		Joins("LEFT JOIN bookmark_tag ON tag.id = bookmark_tag.tag_id").
		Group("tag.id, tag.name")

	if name != "" {
		query = query.Where("tag.name LIKE ?", "%"+name+"%")
	}

	err := query.Order("tag.name ASC").Find(&tags).Error
	return tags, err
}

// FindByID 根据ID查找标签
func (r *TagRepo) FindByID(id int) (*db.Tag, error) {
	var tag db.Tag
	err := lib.DB.First(&tag, id).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Update 更新标签
func (r *TagRepo) Update(tag *db.Tag) error {
	return lib.DB.Updates(tag).Error
}

// Delete 删除标签
func (r *TagRepo) Delete(id int) error {
	return lib.DB.Delete(&db.Tag{}, id).Error
}
