package requests

type SaveProfileRequest struct {
	Login       string `json:"login" validate:"required"`
	Name        string `json:"name,omitempty" validate:"required"`
	Email       string `json:"email,omitempty"`
	Photo       string `json:"photo,omitempty"`
	PhotoTitle  string `json:"photo_filename,omitempty"`
	Description string `json:"description,omitempty"`
	Role        string `json:"role" validate:"required"`
}

type UpdateProfileRequest struct {
	Login       string `json:"login,omitempty"`
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Photo       string `json:"photo,omitempty"`
	PhotoTitle  string `json:"photo_filename,omitempty"`
	Description string `json:"description,omitempty"`
}

type GetByWalletAddressRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
}
