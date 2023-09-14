package enum

//go:generate enumer -type RoleType -json
type RoleType int // defined type
const (
	RoleUnknown RoleType = iota - 1
	RoleUser
	RoleAdmin
)
