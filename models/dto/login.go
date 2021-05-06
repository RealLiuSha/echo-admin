package dto

type Login struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
	CaptchaID   string `json:"captcha_id" validate:"required"`
	CaptchaCode string `json:"captcha_code" validate:"required"`
}
