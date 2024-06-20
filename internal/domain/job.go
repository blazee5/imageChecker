package domain

type CreateJobRequest struct {
	Name      string `json:"name" binding:"required"`
	Image     string `json:"image" binding:"required"`
	IsPrivate bool   `json:"is_private"`
	Username  string `json:"username" binding:"required_if=IsPrivate true"`
	Password  string `json:"password" binding:"required_if=IsPrivate true"`
}
