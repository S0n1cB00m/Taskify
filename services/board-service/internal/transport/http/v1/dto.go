package v1

type CreateBoardRequest struct {
	Title       string `json:"title" example:"Important thing"`
	Description string `json:"description" example:"This is my board's description"`
	Owner       int64  `json:"owner" example:"1"`
}

type UpdateBoardRequest struct {
	Title       *string `json:"title" example:"Important thing"` // Если поля нет в JSON, будет nil
	Description *string `json:"description" example:"This is my board's description"`
}

type ErrBoardNotFoundResponse struct {
	Error string `json:"error" example:"board not found"`
}
