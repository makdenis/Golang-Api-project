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

type PostDetails struct{

	Post Post `json:"post"`
}

type PostDetails2 struct{
	Author User `json:"author,omitempty"`
	Post Post `json:"post"`

}
type PostDetails3 struct{
	Thread Thread `json:"thread"`
	Post Post `json:"post"`

}
type PostDetails6 struct{
	Author User `json:"author,omitempty"`
	Forum Forum `json:"forum"`
	Post Post `json:"post"`

}
type PostDetails4 struct{
	Author User `json:"author,omitempty"`
	Thread Thread `json:"thread"`
	Post Post `json:"post"`

}
type PostDetails5 struct{
	Forum Forum `json:"forum"`
	Post Post `json:"post"`

}
type PostDetails7 struct{

	Forum Forum `json:"forum"`
	Thread Thread `json:"thread"`
	Post Post `json:"post"`

}
type PostDetails8 struct{
	Author User `json:"author,omitempty"`
	Forum Forum `json:"forum"`
	Thread Thread `json:"thread"`
	Post Post `json:"post"`

}