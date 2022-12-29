package auth

type GetUserResponse struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	ProfilePic  string `json:"profile_pic"`
}
