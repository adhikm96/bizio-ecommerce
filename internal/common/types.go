package common

type ReviewCreateDto struct {
	UserID  uint   `json:"user_id"`
	Rating  uint   `json:"rating"`
	Comment string `json:"comment"`
}

type ReviewListDto struct {
	UserID    uint   `json:"user_id"`
	ProductID uint   `json:"product_id"`
	Rating    uint   `json:"rating"`
	Comment   string `json:"comment"`
}
