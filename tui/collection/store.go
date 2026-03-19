package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func baseDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".config", "blink", "collections")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// Save menyimpan collection ke disk.
func Save(name string, data Collection) error {
	dir, err := baseDir()
	if err != nil {
		return fmt.Errorf("base dir: %w", err)
	}

	path := filepath.Join(dir, name+".json")
	tmpPath := path + ".tmp"

	bytes, err := json.MarshalIndent(data, "", "  ")
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
