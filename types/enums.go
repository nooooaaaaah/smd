package types

// enum for roles of users
type Role int

const (
	Admin Role = iota
	Owner
	Regular
	Developer
)
