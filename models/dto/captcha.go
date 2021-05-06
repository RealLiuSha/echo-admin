package dto

type CaptchaVerify struct {
	ID   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}
