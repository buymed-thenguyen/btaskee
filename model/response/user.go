package response

import "time"

type Token struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expire"`
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}
