package authbundle

import core "github.com/rubenscholle/venus/corebundle"

type AuthUser struct {
	core.Model
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email_name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsSysAdmin bool   `json:"-"`
}
