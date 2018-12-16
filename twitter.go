package main

type twitter interface {
	Self() (twitterProfile, error)
	GetFollowerIDs() ([]int64, error)
	SendDM(toID, message string) error
	GetUserByID(id int64) (twitterProfile, error)
}

type twitterProfile struct {
	idStr      string
	screenName string
}
