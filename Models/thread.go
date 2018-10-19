package Models

type Thread struct {
	Author string `json:"author"`

	Created string `json:"created,omitempty"`

	Forum string `json:"forum,omitempty"`

	ID int `json:"id"`

	Message string `json:"message"`

	Slug string `json:"slug,omitempty"`

	Title string `json:"title"`

	Votes int `json:"votes,omitempty"`
}
