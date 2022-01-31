package types

type (
	AuthRequest struct {
		Email    string `json:"email" form:"email" validate:"required"`
		Password string `json:"password" form:"password" validate:"required"`
		Code     string `json:"code" form:"code" validate:"required"`
	}

	LoginLogoutRequest struct {
		Code string `query:"code" validate:"required"`
	}

	PasswordRecoverRequest struct {
		Email string `json:"email" form:"email" validate:"required"`
	}

	PasswordChangeRequest struct {
		Password        string `json:"password" form:"password" validate:"required,eqfield=ConfirmPassword"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required"`
		Code            string `json:"code" form:"code" validate:"required"`
	}

	UserCreateRequest struct {
		Name            string `json:"name" form:"name" validate:"required"`
		Email           string `json:"email" form:"email" validate:"required,email"`
		Password        string `json:"password" form:"password" validate:"required,eqfield=ConfirmPassword"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required"`
		Gender          string `json:"gender" form:"gender" validate:"required"`
		Agreement       bool   `json:"agreement" form:"agreement" validate:"required"`
	}

	UserVerificationRequest struct {
		Token string `query:"token" validate:"required"`
	}

	ApplicationCreateRequest struct {
		Application string `json:"application" validate:"required"`
		Domain      string `json:"domain" validate:"required"`
		RedirectUrl string `json:"redirect_url" validate:"required,url"`
	}
)
