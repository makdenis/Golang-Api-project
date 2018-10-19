package Models

type Post struct {
	Author string `json:"author,omitempty"`

	Created string `json:"created,omitempty"`

	Forum string `json:"forum,omitempty"`

	ID int `json:"id"`

	IsEdited bool `json:"isEdited,omitempty"`

	Message string `json:"message,omitempty"`

	Parent int `json:"parent,omitempty"`

	Thread int `json:"thread,omitempty"`

	//Path []uint64 `json:",omitempty"`
}
