package auth

type SignInPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignOutPayload struct {
	RefreshToken string `json:"refreshToken"`
}

type SignUpPayload struct {
	SignInPayload
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
