package account

type ChangePasswordPayload struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

type UpdateAccountPayload struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
