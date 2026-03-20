package utils

import (
	"fmt"
	"go-httpix-cli/entity"
	"log"
	"os"
	"strings"
)

var logFile, _ = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
var Logger = log.New(logFile, "", log.LstdFlags)

func TableRowsToEnv(name string, rows []entity.EnvTableRow) entity.Env {
	vars := make(map[string]string, len(rows))
	for _, row := range rows {
		k := strings.TrimSpace(row.Key.Value())
		v := strings.TrimSpace(row.Value.Value())
		if k != "" {
			vars[k] = v
		}
	}
	return entity.Env{Name: name, Vars: vars}
}

// Flatten mengubah nested Collection menjadi flat []TreeNode
// hanya node yang visible (parent-nya expanded) yang dimasukkan
func Flatten(collections []entity.Collection) []entity.TreeNode {
	var result []entity.TreeNode
	for _, c := range collections {
		flattenNode(c, 0, &result)
	}
	return result
}

func flattenNode(c entity.Collection, depth int, result *[]entity.TreeNode) {
	node := entity.TreeNode{
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
		*result = append(*result, entity.TreeNode{
			ID:       req.ID,
			Name:     req.Name,
			Depth:    depth + 1,
			IsFolder: false,
			Data:     &req,
		})
	}
}

func BuildURL(rawURL, params string) (string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = fmt.Sprintf("https://%s", rawURL)
	}
	if queryParam := strings.TrimSpace(params); queryParam != "" {
		stringSeparator := "?"
		if strings.Contains(rawURL, "?") {
			stringSeparator = "&"
		}
		rawURL += stringSeparator + strings.ReplaceAll(queryParam, "\n", "&")
	}
	return rawURL, nil
}

func ApplyHeaders(rawString string) map[string]string {
	mappedHeaders := make(map[string]string)

	for _, rawRow := range strings.Split(rawString, "\n") {
		rawRow = strings.TrimSpace(rawRow)
		if rawRow == "" {
			continue
		}

		splittedPart := strings.SplitN(rawRow, "=", 2)
		if len(splittedPart) != 2 {
			continue // skip invalid format
		}

		key := strings.TrimSpace(splittedPart[0])
		value := strings.TrimSpace(splittedPart[1])

		mappedHeaders[key] = value
	}

	return mappedHeaders
}
