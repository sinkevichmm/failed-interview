package models

import "time"

type File struct {
	ID         string
	Name       string
	Extension  string
	DateUpload time.Time
	File       []byte
}
