package owner

type Original struct {
	ID               string
	Email            string
	IsAuthedStreamer bool
	AccessToken      string
	RefreshToken     string
	Address          string
	Channel
}

type Channel struct {
	Name        string
	Description string
	Image       string
	URL         string
}

type User struct {
	Address  string
	Nickname string
}
