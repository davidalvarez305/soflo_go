package types

type User struct {
	Username string `json:"username"`
	Password []byte `json:"password"`
}

type UserRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
