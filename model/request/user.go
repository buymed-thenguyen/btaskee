package request

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}
