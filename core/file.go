package core

import (
	"encoding/json"
	"fmt"
	"go-httpix-cli/tui/collection"
	"os"
	"path/filepath"
)

func baseDir() (string, error) {
	home, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, "collections")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func SaveFile(collectionName string, collection collection.Collection) error {
	dir, err := baseDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, collectionName+".json")
	tmpPath := path + ".tmp"
	bytes, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	if err := os.WriteFile(tmpPath, bytes, 0644); err != nil {
		return fmt.Errorf("write tmp: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	return nil
}
