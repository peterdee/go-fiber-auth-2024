package account

type ChangePasswordPayload struct {
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}
