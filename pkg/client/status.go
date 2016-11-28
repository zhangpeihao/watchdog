package client

type Status struct {
	URL string
}

func GetStatus() []*Status {
	return []*Status{
		&Status{
			URL: "http://xxxx",
		},
		&Status{
			URL: "http://yyyy",
		},
	}
}
