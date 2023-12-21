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
	ID          string
	Name        string
	Size        int64
	ContentType string
	Location    string
	UploadDate  time.Time
	OwnerID     string
}

type User struct {
	ID        string
	Username  string
	Password  string
	Email     string
	Role      Role
	CreatedAt time.Time
}

type AccessControlList []User

type Session struct {
	ID     string
	UserID string
	Token  AuthToken
}

type ActiveSessions []Session

type AuthToken struct {
	Token     string
	ExpiresAt time.Time
}

type ApiResponse struct {
	Success bool
	Message string
	Data    interface{}
}

type ApiError struct {
	Code    int
	Key     string
	Message string
}

type RolePrivileges map[Role]Privileges
