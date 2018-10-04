package db

import "database/sql"

type Host struct {
	ID             int
	IP             string
	HostName       string
	IsInit         bool
	PresentVersion Version
	Spec           Spec
	Files          []File
}

type Spec struct {
	ID    int
	Cpu   string
	Mem   string
	Disk  string
	Hosts []Host
}

type File struct {
	ID             int
	Name           string
	Path           string
	Version        string
	PresentVersion Version
}

type Version struct {
	Version sql.NullString
}
