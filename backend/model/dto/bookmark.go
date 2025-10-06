package dto

// BookmarkListRequest 书签列表请求
type BookmarkListRequest struct {
	Keyword  string `form:"keyword" json:"keyword"`                    // 内容查询关键字, 查询范围：url、title、excerpt、content
	Tags     string `form:"tags" json:"tags"`                          // tag列表，使用英文的逗号分隔多个tag，tag name全匹配
	Page     int    `form:"page" json:"page" binding:"required,min=1"` // 页码
	PageSize int    `form:"page_size" json:"page_size"`                // 每页记录数，默认10
}

// BookmarkListItem 书签列表项
type BookmarkListItem struct {
	ID        int       `json:"id"`
	URL       string    `json:"url"`            // 内容的原文地址
	Title     string    `json:"title"`          // 标题
	Excerpt   string    `json:"excerpt"`        // 摘录，可以是自定的概述
	Author    string    `json:"author"`         // 作者
	IsArchive bool      `json:"is_archive"`     // 是否已归档
	CreatedAt int64     `json:"created_at"`     // 创建时间（时间戳）
	UpdatedAt int64     `json:"updated_at"`     // 最后更新时间（时间戳）
	Tags      []TagItem `json:"tags,omitempty"` // 标签列表
}

// BookmarkListResponse 书签列表响应
type BookmarkListResponse struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	Data PageData `json:"data"`
}

// CreateBookmarkRequest 创建书签请求
type CreateBookmarkRequest struct {
	URL           string    `json:"url" binding:"required"` // 书签原文地址
	Title         string    `json:"title"`                  // 书签标题
	Excerpt       string    `json:"excerpt"`                // 简述
	Tags          []TagItem `json:"tags"`                   // 标签列表
	CreateArchive bool      `json:"create_archive"`         // 是否创建归档
}

// UpdateBookmarkRequest 编辑书签请求
type UpdateBookmarkRequest struct {
	ID            int       `json:"id" binding:"required"`  // 书签ID
	URL           string    `json:"url" binding:"required"` // 书签原文地址
	Title         string    `json:"title"`                  // 书签标题
	Excerpt       string    `json:"excerpt" `               // 简述
	Author        string    `json:"author" `                // 作者
	Tags          []TagItem `json:"tags"`                   // 标签列表
	CreateArchive bool      `json:"create_archive"`         // 是否创建归档
}

// DeleteBookmarkRequest 删除书签请求（批量删除）
type DeleteBookmarkRequest []int

// BookmarkContentResponse 书签内容响应数据
type BookmarkContentData struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`        // 书签原文地址
	Title     string `json:"title"`      // 标题
	HTML      string `json:"html"`       // 内容
	CreatedAt int64  `json:"created_at"` // 创建时间（时间戳）
	UpdatedAt int64  `json:"update_at"`  // 最近更新时间（时间戳）
}

// BookmarkContentResponse 书签内容响应
type BookmarkContentResponse struct {
	Code int                 `json:"code"`
	Msg  string              `json:"msg"`
	Data BookmarkContentData `json:"data"`
}

// ImportProgressEvent SSE 导入进度事件
type ImportProgressEvent struct {
	Type    string `json:"type"`    // progress, success, error, complete
	Message string `json:"message"` // 消息内容
	Current int    `json:"current"` // 当前处理数量
	Total   int    `json:"total"`   // 总数量
	URL     string `json:"url"`     // 当前处理的 URL
}
