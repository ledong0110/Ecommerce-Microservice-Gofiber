package models

type EmailForm struct {
	Email string `json:"email"`
}

type NewPasswordForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePasswordForm struct {
	Username string `json:"username"`
	OldPassword string `json:"old_pwd"`
	NewPassword string `json:"new_pwd"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterForm struct {
	Fullname string `json:"fullname"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Picture  string `json:"picture"`
	Password string `json:"password"`
}

type OTPForm struct {
	Email string `json:"email"`
	OTP   string `json:"otp"`
}
