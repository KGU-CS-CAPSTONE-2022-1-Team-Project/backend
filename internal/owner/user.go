package owner

type User struct {
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
