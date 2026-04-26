package models

//User : id, name
type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Password   string `json:"-"`
	Email      string `json:"email"`
	ProfilePic string `json:"profile_pic"`
}
