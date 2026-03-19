package collection

type TreeNode struct {
	ID       string
	Name     string
	Depth    int      // kedalaman indent (0 = root)
	IsFolder bool     // true = Collection, false = Request
	Expanded bool     // hanya relevan kalau IsFolder
	Data     *Request // nil kalau IsFolder
}

// Flatten mengubah nested Collection menjadi flat []TreeNode
// hanya node yang visible (parent-nya expanded) yang dimasukkan
func Flatten(collections []Collection) []TreeNode {
	var result []TreeNode
	for _, c := range collections {
		flattenNode(c, 0, &result)
	}
	return result
}

func flattenNode(c Collection, depth int, result *[]TreeNode) {
	node := TreeNode{
		ID:       c.ID,
		Name:     c.Name,
		Depth:    depth,
		IsFolder: true,
		Expanded: c.Expanded,
	}
	*result = append(*result, node)

	// kalau tidak expanded, anak-anaknya tidak dimasukkan
	if !c.Expanded {
		return
	}

	// masukkan sub-collection
	for _, child := range c.Children {
		flattenNode(child, depth+1, result)
	}

	// masukkan request
	for _, req := range c.Requests {
		*result = append(*result, TreeNode{
			ID:       req.ID,
			Name:     req.Name,
			Depth:    depth + 1,
			IsFolder: false,
			Data:     &req,
		})
	}
}
