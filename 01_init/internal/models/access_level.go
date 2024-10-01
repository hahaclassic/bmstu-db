package models

type AccessLevel int

const (
	Owner AccessLevel = iota + 1
	Other
)
