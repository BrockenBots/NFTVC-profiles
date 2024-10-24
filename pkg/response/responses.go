package response

type SaveProfileResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Error string `json:"error" validate:"required"`
}

type GetMeResponse struct {
	Id          string `json:"profile_id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Photo       string `json:"photo"`
	Description string `json:"description"`
	AccountId   string `json:"account_id"`
}

type UpdateProfileResponse struct {
	Id          string `json:"profile_id"`
	Login       string `json:"login"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Photo       string `json:"photo"`
	Description string `json:"description"`
	AccountId   string `json:"account_id"`
}
