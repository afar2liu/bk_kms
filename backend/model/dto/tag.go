package dto

// TagListRequest tag列表请求
type TagListRequest struct {
	Name string `form:"name" json:"name"` // 按名称查询tag
}

// TagListItem tag列表项
type TagListItem struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`  // 标签名称
	Count int    `json:"count"` // 绑定标签的文档数量
}

// TagListResponse tag列表响应
type TagListResponse struct {
	Code int           `json:"code"`
	Msg  string        `json:"msg"`
	Data []TagListItem `json:"data"`
}

// UpdateTagRequest tag重命名请求
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"` // 新的tag名称
}
