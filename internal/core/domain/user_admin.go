package domain

type UserFilter struct {
	Role      string
	IsActive  *bool
	IsBlocked *bool
	Search    string
	Page      int
	Limit     int
}

type UserPatch struct {
	Email         *string
	PasswordHash  *string
	FullName      *string
	PhoneProvided bool
	Phone         *string
	Role          *Role
	IsActive      *bool
	IsBlocked     *bool
}
