package outbound

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type CollectionIndex struct {
	Version int      `json:"version"`
	Order   []string `json:"order"`
}

func indexPath() (string, error) {
	dir, err := baseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "index.json"), nil
}

func LoadIndex() (CollectionIndex, error) {
	path, err := indexPath()
	if err != nil {
		return CollectionIndex{Version: 1, Order: []string{}}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return CollectionIndex{Version: 1, Order: []string{}}, nil
	}
	if err != nil {
		return CollectionIndex{}, err
	}

	var idx CollectionIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return CollectionIndex{}, err
	}
	return idx, nil
}

func SaveIndex(idx CollectionIndex) error {
	path, err := indexPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func AppendToIndex(id string) error {
	idx, err := LoadIndex()
	if err != nil {
		return err
	}
	for _, existing := range idx.Order {
		if existing == id {
			return nil // sudah ada
		}
	}
	idx.Order = append(idx.Order, id)
	return SaveIndex(idx)
}

func RemoveFromIndex(id string) error {
	idx, err := LoadIndex()
	if err != nil {
		return err
	}
	newOrder := make([]string, 0, len(idx.Order))
	for _, existing := range idx.Order {
		if existing != id {
			newOrder = append(newOrder, existing)
		}
	}
	idx.Order = newOrder
	return SaveIndex(idx)
}
