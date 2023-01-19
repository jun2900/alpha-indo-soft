package entity

type CreateArticleReq struct {
	Author string `json:"author"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
