package dto

type CustomClaimDto struct {
	Key   string `json:"key" binding:"required,max=20"`
	Value string `json:"value" binding:"required,max=10000"`
}

type CustomClaimCreateDto = CustomClaimDto
