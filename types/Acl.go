package types

import (
	"errors"
)

func (acl *AccessControlList) AddUser(user User) error {
	for _, existingUser := range *acl {
		if existingUser == user {
			return errors.New("user already exists in the list")
		}
	}
	*acl = append(*acl, user)
	return nil
}

func (acl *AccessControlList) RemoveUser(user User) error {
	for i, existingUser := range *acl {
		if existingUser == user {
			*acl = append((*acl)[:i], (*acl)[i+1:]...)
			return nil
		}
	}
	return errors.New("user not found in the list")
}

func (acl *AccessControlList) HasUser(user User) bool {
	for _, existingUser := range *acl {
		if existingUser == user {
			return true
		}
	}
	return false
}

func (acl *AccessControlList) HasUserWithId(id string) bool {
	for _, existingUser := range *acl {
		if existingUser.ID == id {
			return true
		}
	}
	return false
}
