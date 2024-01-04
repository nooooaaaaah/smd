package types

import "time"

type Privileges struct {
	Read              bool
	Write             bool
	Delete            bool
	CreateDirectories bool
	AddUsers          bool
}

type Directory struct {
	ID                string
	Name              string
	OwnerID           string
	ParentDirectoryID string
	FileIDs           []string
	SubdirectoryIDs   []string
}

type File struct {
	UploadDate  time.Time
	ID          string
	Name        string
	ContentType string
	Location    string
	OwnerID     string
	Size        int64
}

type User struct {
	CreatedAt time.Time
	ID        string
	Username  string
	Password  string
	Email     string
	Role      Role
}

type AccessControlList []User

type Session struct {
	Token  AuthToken
	ID     string
	UserID string
}

type ActiveSessions []Session

type AuthToken struct {
	ExpiresAt time.Time
	Token     string
}

type ApiResponse struct {
	Data    interface{}
	Message string
	Success bool
}

type ApiError struct {
	Key     string
	Message string
	Code    int
}

type RolePrivileges map[Role]Privileges
