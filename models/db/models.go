package db

import "database/sql"

type Host struct {
	ID             int
	IP             string
	HostName       string
	IsInit         bool
	SpecID         int
	PresentVersion Version
	Spec           Spec
	Files          []File
}

type Spec struct {
	ID    int
	Cpu   int
	Mem   int
	Disk  int
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
