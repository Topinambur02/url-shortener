package dto

type CreateURLDto struct {
	OriginalURL string `json:"original_url" validate:"required,url" example:"https://finance.ozon.ru"`
}
