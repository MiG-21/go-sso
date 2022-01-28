package types

type (
	AuthRequest struct {
		Email    string `query:"email" validate:"required"`
		Password string `query:"password" validate:"required"`
	}

	UserCreateRequest struct {
		Name            string `json:"name" form:"name" validate:"required"`
		Email           string `json:"email" form:"email" validate:"required,email"`
		Password        string `json:"password" form:"password" validate:"required,eqfield=ConfirmPassword"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" validate:"required"`
		Gender          string `json:"gender" form:"gender" validate:"required"`
	}
)
