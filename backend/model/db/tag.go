package db

// Tag tag表
type Tag struct {
	ID   int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"column:name;type:varchar(250);not null;uniqueIndex:tag_name_UNIQUE" json:"name"`

	// 关联关系
	Bookmarks []Bookmark `gorm:"many2many:bookmark_tag;foreignKey:ID;joinForeignKey:tag_id;References:ID;joinReferences:bookmark_id" json:"bookmarks,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tag"
}
