package Models

type User struct {
	Fullname	string	`json:"fullname,omitempty"`
	Nickname	string	`json:"nickname,omitempty"`
	Email  		string	`json:"email,omitempty"`
	About		string 	`json:"about,omitempty"`

}

