package enum

//go:generate go run github.com/alvaroloes/enumer -type RoleType -json
type RoleType int // defined type
const (
	RoleUnknown RoleType = iota - 1
	RoleUser
	RoleAdmin
)
