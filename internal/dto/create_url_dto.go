package dto

type CreateUrlDto struct {
	OriginalUrl string `json:"original_url" validate:"required,url"`
}
