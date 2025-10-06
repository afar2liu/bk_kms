package dto

// Response 通用响应结构
type Response struct {
	Code int         `json:"code"` // 0:成功，>0:失败
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// PageData 分页数据结构
type PageData struct {
	Rows  interface{} `json:"rows"`
	Total int         `json:"total"`
}

// TagItem Tag项
type TagItem struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
}
