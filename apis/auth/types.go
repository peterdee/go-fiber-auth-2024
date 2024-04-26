package auth

type RefreshTokensPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

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
