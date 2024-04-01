package auth

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpPayload struct {
	SignInPayload
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
