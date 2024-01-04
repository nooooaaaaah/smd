package types

import (
	"database/sql"
)

type Database interface {
	CreateDb() error
	Connect() error
	Disconnect() error
	InsertUser(u User) error
	InsertFile(f File) error
	InsertDirectory(d Directory) error
	GetUser(username string) (User, error)
	GetFile(filename string) (File, error)
	GetDirectory(directoryName string) (Directory, error)
	GetAllFiles() ([]File, error)
	GetAllDirectories() ([]Directory, error)
	GetAllUsers() ([]User, error)
	UpdateUser(u User) error
	UpdateFile(f File) error
	UpdateDirectory(d Directory) error
	DeleteUser(username string) error
	DeleteFile(filename string) error
	DeleteDirectory(directoryName string) error
}

type database struct {
	db *sql.DB
}

func NewDatabase() Database {
	return &database{}
}

func (d *database) CreateDb() error {
	d.Connect()
	defer d.Disconnect()
	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, username TEXT, password TEXT, email TEXT, role TEXT, created_at TEXT)")
	if err != nil {
		return err
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS files (id TEXT PRIMARY KEY, name TEXT, size TEXT, content_type TEXT, location TEXT, upload_date TEXT, owner_id TEXT)")
	if err != nil {
		return err
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS directories (id TEXT PRIMARY KEY, name TEXT, owner_id TEXT, parent_directory_id TEXT)")
	if err != nil {
		return err
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS access_control_lists (id TEXT PRIMARY KEY, user_id TEXT, file_id TEXT, directory_id TEXT)")
	if err != nil {
		return err
	}
	_, err = d.db.Exec("CREATE TABLE IF NOT EXISTS active_sessions (id TEXT PRIMARY KEY, user_id INTEGER, token TEXT, expires_at TEXT)")
	if err != nil {
		return err
	}

	return nil
}

func (d *database) Connect() error {
	db, err := sql.Open("sqlite3", "./Smd.db")
	if err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *database) Disconnect() error {
	return d.db.Close()
}

// DeleteDirectory in database
func (d *database) DeleteDirectory(directoryName string) error {
	_, err := d.db.Exec("DELETE FROM directories WHERE name = ?", directoryName)
	if err != nil {
		return err
	}
	return nil
}

// DeleteFile in database
func (d *database) DeleteFile(filename string) error {
	_, err := d.db.Exec("DELETE FROM files WHERE name = ?", filename)
	if err != nil {
		return err
	}
	return nil
}

// DeleteUser in database
func (d *database) DeleteUser(username string) error {
	_, err := d.db.Exec("DELETE FROM users WHERE username = ?", username)
	if err != nil {
		return err
	}
	return nil
}

// GetAllDirectories in database
func (d *database) GetAllDirectories() ([]Directory, error) {
	rows, err := d.db.Query("SELECT * FROM directories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var directories []Directory
	for rows.Next() {
		var directory Directory
		err := rows.Scan(&directory.ID, &directory.Name, &directory.OwnerID, &directory.ParentDirectoryID)
		if err != nil {
			return nil, err
		}
		directories = append(directories, directory)
	}

	return directories, nil
}

// GetAllFiles in database
func (d *database) GetAllFiles() ([]File, error) {
	rows, err := d.db.Query("SELECT * FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []File
	for rows.Next() {
		var file File
		err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &file.Location, &file.UploadDate, &file.OwnerID)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

// GetAllUsers
func (d *database) GetAllUsers() ([]User, error) {
	rows, err := d.db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// GetDirectory in database
func (d *database) GetDirectory(directoryName string) (Directory, error) {
	row := d.db.QueryRow("SELECT * FROM directories WHERE name = ?", directoryName)
	var directory Directory
	err := row.Scan(&directory.ID, &directory.Name, &directory.OwnerID, &directory.ParentDirectoryID)
	if err != nil {
		return Directory{}, err
	}
	return directory, nil
}

// GetFile in database
func (d *database) GetFile(filename string) (File, error) {
	row := d.db.QueryRow("SELECT * FROM files WHERE name = ?", filename)
	var file File
	err := row.Scan(&file.ID, &file.Name, &file.Size, &file.ContentType, &file.Location, &file.UploadDate, &file.OwnerID)
	if err != nil {
		return File{}, err
	}
	return file, nil
}

// GetUser in database
func (d *database) GetUser(username string) (User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE username = ?", username)
	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// InsertDirectory in database
func (d *database) InsertDirectory(dir Directory) error {
	_, err := d.db.Exec("INSERT INTO directories (id, name, owner_id, parent_directory_id) VALUES (?, ?, ?, ?)", dir.ID, dir.Name, dir.OwnerID, dir.ParentDirectoryID)
	if err != nil {
		return err
	}
	return nil
}

// InsertFile in database
func (d *database) InsertFile(f File) error {
	_, err := d.db.Exec("INSERT INTO files (id, name, size, content_type, location, upload_date, owner_id) VALUES (?, ?, ?, ?, ?, ?, ?)", f.ID, f.Name, f.Size, f.ContentType, f.Location, f.UploadDate, f.OwnerID)
	if err != nil {
		return err
	}
	return nil
}

// InsertUser in database
func (d *database) InsertUser(u User) error {
	_, err := d.db.Exec("INSERT INTO users (id, username, password, email, role, created_at) VALUES (?, ?, ?, ?, ?, ?)", u.ID, u.Username, u.Password, u.Email, u.Role, u.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

// UpdateDirectory in database
func (d *database) UpdateDirectory(dir Directory) error {
	_, err := d.db.Exec("UPDATE directories SET name = ?, owner_id = ?, parent_directory_id = ? WHERE id = ?", dir.Name, dir.OwnerID, dir.ParentDirectoryID, dir.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateFile in database
func (d *database) UpdateFile(f File) error {
	_, err := d.db.Exec("UPDATE files SET name = ?, size = ?, content_type = ?, location = ?, upload_date = ?, owner_id = ? WHERE id = ?", f.Name, f.Size, f.ContentType, f.Location, f.UploadDate, f.OwnerID, f.ID)
	if err != nil {
		return err
	}
	return nil
}

// UpdateUser in database
func (d *database) UpdateUser(u User) error {
	_, err := d.db.Exec("UPDATE users SET username = ?, password = ?, email = ?, role = ?, created_at = ? WHERE id = ?", u.Username, u.Password, u.Email, u.Role, u.CreatedAt, u.ID)
	if err != nil {
		return err
	}
	return nil
}
