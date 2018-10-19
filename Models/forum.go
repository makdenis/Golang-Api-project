package Models

type Forum struct {
	User string `json:"user,omitempty"`
	Title string `json:"title,omitempty"`
	Slug string `json:"slug,omitempty"`
	Posts int `json:"posts,omitempty"`
	Threads int `json:"threads,omitempty"`
}

