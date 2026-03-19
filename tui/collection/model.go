package collection

type Collection struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Children []Collection `json:"children"` // rekursif
	Requests []Request    `json:"requests"`
	Expanded bool         `json:"-"`
}

type Request struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Method  string `json:"method"`
	URL     string `json:"url"`
	Headers string `json:"headers"`
	Body    string `json:"body"`
	Params  string `json:"params"`
}
