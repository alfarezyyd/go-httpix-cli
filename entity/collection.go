package entity

type Collection struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	Children []Collection `json:"children"` // rekursif
	Requests []Request    `json:"requests"`
	Expanded bool         `json:"-"`
}

type TreeNode struct {
	ID       string
	Name     string
	Depth    int      // kedalaman indent (0 = root)
	IsFolder bool     // true = Collection, false = Request
	Expanded bool     // hanya relevan kalau IsFolder
	Data     *Request // nil kalau IsFolder
}
