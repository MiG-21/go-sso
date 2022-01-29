package types

type (
	UserTokenResponse struct {
		Token string `json:"token"`
	}

	UserCreateResponse struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	UserInfoResponse struct {
		Id       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Gender   string `json:"gender"`
		Data     string `json:"data"`
		Created  int64  `json:"created"`
		Updated  int64  `json:"updated"`
		Active   bool   `json:"active"`
		Locked   bool   `json:"locked"`
		LockedTo int64  `json:"locked_to"`
	}
)
