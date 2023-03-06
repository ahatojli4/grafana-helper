package entities

type SearchDashboard struct {
	Uid   string        `json:"uid"`
	Title string        `json:"title"`
	Url   string        `json:"url"`
	Type  string        `json:"type"`
	Tags  []interface{} `json:"tags"`
}
