package models

type Role string

const (
	RoleAdmin      Role = "ADMIN"
	RoleSuperAdmin Role = "SUPERADMIN"
	RoleUser       Role = "USER"
)
