package requset

type UserSignUp struct {
	Telephone string
	Username  string
	Password  string
}

type UserSignIn struct {
	Telephone string
	Password  string
}
