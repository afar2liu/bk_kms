package db

import "time"

// Bookmark 书签表
type Bookmark struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	URL       string    `gorm:"column:url;type:text;not null;uniqueIndex:bookmark_url_UNIQUE,length:255;comment:网址地址" json:"url"`
	Title     string    `gorm:"column:title;type:text;not null;comment:网址标题" json:"title"`
	Excerpt   string    `gorm:"column:excerpt;type:text;not null;comment:网站内容节选" json:"excerpt"`
	Author    string    `gorm:"column:author;type:text;not null;comment:作者" json:"author"`
	Content   string    `gorm:"column:content;type:mediumtext;not null;comment:网站文本内容(去掉html标签)" json:"content"`
	HTML      string    `gorm:"column:html;type:mediumtext;not null;comment:网站原始内容" json:"html"`
	IsArchive bool      `gorm:"column:is_archive;type:tinyint(1);not null;default:0;comment:是否已存档,0:否，1:是" json:"is_archive"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime;index:idx_created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime;index:idx_modified_at" json:"updated_at"`

	// 关联关系
	Tags []Tag `gorm:"many2many:bookmark_tag;foreignKey:ID;joinForeignKey:bookmark_id;References:ID;joinReferences:tag_id" json:"tags,omitempty"`
}

// TableName 指定表名
func (Bookmark) TableName() string {
	return "bookmark"
}
