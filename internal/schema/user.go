package schema

type NewUserRequest struct {
	OrgID       int64  `json:"organization_id" validate:"-"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

type UserCompact struct {
	ID          int64  `json:"id"`
	DisplayName string `json:"displayName"`
	Email       string `json:"username"`
	CreatedAt   string `json:"createdAt"`
}

// CompactUserListResponse list of users but compact form
type CompactUserListResponse struct {
	Total int           `json:"total"`
	Users []UserCompact `json:"users"`
}

type UpdateUserRequest struct {
	ID          int32  `json:"id" validate:"-"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	DisplayName string `json:"display_name" validate:"required"`
	Phone       string `json:"phone" validate:"e164"`
	OrgID       int32  `json:"org_id" validate:"required"`
	CountryIso  string `json:"country_iso" validate:"iso3166_1_alpha2"`
	City        string `json:"city" validate:"required"`
	Address     string `json:"address" validate:"-"`
}
