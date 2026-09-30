package dto

type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required,min=3,max=320"`

	Password string `json:"password" validate:"required,min=1,max=1024"`
}
