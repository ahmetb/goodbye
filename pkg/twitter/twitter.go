package twitter

type Twitter interface {
	Self() (TwitterProfile, error)
	GetFollowerIDs() ([]int64, error)
	SendDM(toID, message string) error
	GetUserByID(id int64) (TwitterProfile, error)
}

type TwitterProfile struct {
	IDStr      string
	ScreenName string
}
