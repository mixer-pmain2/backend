package types

type ChangePassword struct {
	UserId      int64  `json:"userId"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
