package models

type AccessLevel string

const (
	AccessLevelWorker  AccessLevel = "worker"
	AccessLevelManager AccessLevel = "manager"
	AccessLevelAdmin   AccessLevel = "admin"
	AccessLevelOwner   AccessLevel = "owner"
)
