package common

type UserReg struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
