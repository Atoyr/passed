package Models

type Signin struct {
	Email    string `json:"email"`
	Key      string `json:"key"`
	Password string `json:"password"`
}