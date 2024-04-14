package entity

type Content struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

type Banner struct {
	Id         *int64   `json:"id"`
	Content    *Content `json:"content"`
	Feature_id *int64   `json:"feature_id"`
	Tag_id     *int64   `json:"tag_id"`
	Is_active  *bool    `json:"is_active"`
	Tag_ids    []int64  `json:"tag_ids"`
	Created_at string   `json:"created_at"`
	Updated_at string   `json:"updated_at"`
}
