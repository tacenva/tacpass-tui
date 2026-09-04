package sourceoftruth

type Model struct {
	Width  int
	Height int

	// Application data
	Vault Vault

	// UI state
	Cursor int
}

type Vault struct {
	Name string

	Collections []Collection
}

type Collection struct {
	ID   string
	Name string

	Items []Item
}

type Item struct {
	ID       string
	Name     string
	Username string
	Password string
	URL      string
}

func New() Model {
	return Model{
		Vault: Vault{
			Name: "Tacpass Vault",

			Collections: []Collection{
				{
					ID:   "personal",
					Name: "Personal",
					Items: []Item{
						{
							ID:       "github",
							Name:     "GitHub",
							Username: "user@example.com",
							Password: "secret",
							URL:      "https://github.com",
						},
						{
							ID:       "google",
							Name:     "Google",
							Username: "user@example.com",
							Password: "secret",
							URL:      "https://google.com",
						},
					},
				},
				{
					ID:   "work",
					Name: "Work",
					Items: []Item{
						{
							ID:       "aws",
							Name:     "AWS",
							Username: "admin",
							Password: "secret",
							URL:      "https://aws.amazon.com",
						},
					},
				},
			},
		},
	}
}
