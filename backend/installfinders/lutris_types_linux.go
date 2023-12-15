package installfinders

type LutrisGame struct {
	ID        int    `json:"id"`
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	Runner    string `json:"runner"`
	Directory string `json:"directory"`
}
