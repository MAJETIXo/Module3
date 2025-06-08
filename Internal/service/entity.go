package service

type TaskRequest struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}
