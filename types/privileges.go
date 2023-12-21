package types

import (
	"errors"
)

func (p *Privileges) SetPrivileges(read, write, delete, createDirectories, addUsers bool, user User) error {
	if !rolePrivileges[user.Role].AddUsers {
		return errors.New("User does not have privileges to set privileges")
	}
	p.Read = read
	p.Write = write
	p.Delete = delete
	p.CreateDirectories = createDirectories
	p.AddUsers = addUsers
	return nil
}

func (p *Privileges) GetPrivileges(user User) Privileges {
	if p.AddUsers {
		return *p
	}
	return rolePrivileges[user.Role]
}

// Get privileges for a role (Admin, Owner, Regular, Developer)
func (p *Privileges) GetPrivilegesForRole(role Role) Privileges {
	return rolePrivileges[role]
}

// Customize privileges for a role (Admin, Owner, Regular, Developer)
func (p *Privileges) EditPrivileges(user User, newPrivileges *Privileges) error {
	if !rolePrivileges[user.Role].AddUsers {
		return errors.New("User does not have privileges to edit privileges")
	}
	p.Read = newPrivileges.Read
	p.Write = newPrivileges.Write
	p.Delete = newPrivileges.Delete
	p.CreateDirectories = newPrivileges.CreateDirectories
	p.AddUsers = newPrivileges.AddUsers
	return nil
}

// Default Privileges for each role
var rolePrivileges RolePrivileges = map[Role]Privileges{
	Admin: {
		Read:              true,
		Write:             true,
		Delete:            true,
		CreateDirectories: true,
		AddUsers:          true,
	},
	Owner: {
		Read:              true,
		Write:             true,
		Delete:            true,
		CreateDirectories: true,
		AddUsers:          false,
	},
	Regular: {
		Read:              true,
		Write:             false,
		Delete:            false,
		CreateDirectories: false,
		AddUsers:          false,
	},
	Developer: {
		Read:              true,
		Write:             true,
		Delete:            false,
		CreateDirectories: true,
		AddUsers:          false,
	},
}
