package db

// BookmarkTag bookmark与tag关联中间表
type BookmarkTag struct {
	BookmarkID int `gorm:"column:bookmark_id;primaryKey;not null;index:bookmark_tag_bookmark_id_FK" json:"bookmark_id"`
	TagID      int `gorm:"column:tag_id;primaryKey;not null;index:bookmark_tag_tag_id_FK" json:"tag_id"`
}

// TableName 指定表名
func (BookmarkTag) TableName() string {
	return "bookmark_tag"
}
