package domain

type Role string

const (
	RoleClient Role = "client"
	RoleLogist Role = "logist"
	RoleWorker Role = "worker"
	RoleAdmin  Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsValid() bool {
	switch r {
	case RoleClient, RoleLogist, RoleWorker, RoleAdmin:
		return true
	default:
		return false
	}
}
