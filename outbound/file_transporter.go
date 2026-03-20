package outbound

import (
	"encoding/json"
	"fmt"
	"go-httpix-cli/entity"
	"os"
	"path/filepath"
	"strings"
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

func SaveCollection(collectionName string, collection entity.Collection) error {
	dir, err := baseDir()
	if err != nil {
		return err
	}
	path := filepath.Join(dir, collection.ID+".json")
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

	// 2. daftarkan ke index.json
	if err := AppendToIndex(collection.ID); err != nil {
		// kalau index gagal, hapus file yang sudah dibuat
		// agar tidak ada orphan file
		dir, _ := baseDir()
		os.Remove(filepath.Join(dir, collection.ID+".json"))
		return fmt.Errorf("append to index: %w", err)
	}

	return nil
}

func SaveRequestToCollection(req entity.Request, targetFolderID string) error {
	collections, err := LoadAllCollections()
	if err != nil {
		return fmt.Errorf("load collections: %w", err)
	}

	// targetFolderID kosong = simpan di root collection pertama
	if targetFolderID == "" {
		if len(collections) == 0 {
			return fmt.Errorf("no collection exists — create one first")
		}
		root := collections[0]
		root.Requests = append(root.Requests, req)
		return SaveCollection(root.Name, root)
	}

	// cari root collection mana yang mengandung targetFolderID
	for _, root := range collections {
		if !containsID(root, targetFolderID) {
			continue // bukan di sini
		}

		// ketemu root collection-nya
		// targetFolderID == root.ID → sisipkan di root langsung
		if root.ID == targetFolderID {
			root.Requests = append(root.Requests, req)
			return SaveCollection(root.Name, root)
		}

		// targetFolderID ada di dalam nested children
		updated, found := insertRequestInTree(root, targetFolderID, req)
		if !found {
			return fmt.Errorf("folder %s not found", targetFolderID)
		}
		return SaveCollection(updated.Name, updated)
	}

	return fmt.Errorf("folder %s not found in any collection", targetFolderID)
}

// insertRequestInTree menelusuri children secara rekursif
func insertRequestInTree(c entity.Collection, targetID string, req entity.Request) (entity.Collection, bool) {
	for i, child := range c.Children {
		if child.ID == targetID {
			c.Children[i].Requests = append(c.Children[i].Requests, req)
			return c, true
		}

		updated, found := insertRequestInTree(child, targetID, req)
		if found {
			c.Children[i] = updated
			return c, true
		}
	}
	return c, false
}

// containsID cek apakah collection ini atau anak-anaknya mengandung ID
func containsID(c entity.Collection, id string) bool {
	if c.ID == id {
		return true
	}
	for _, r := range c.Requests {
		if r.ID == id {
			return true
		}
	}
	for _, child := range c.Children {
		if containsID(child, id) {
			return true
		}
	}
	return false
}

// LoadCollection membaca satu file collection dari disk
// LoadCollection membaca satu file {id}.json dari disk
func LoadCollection(id string) (entity.Collection, error) {
	dir, err := baseDir()
	if err != nil {
		return entity.Collection{}, err
	}

	data, err := os.ReadFile(filepath.Join(dir, id+".json"))
	if err != nil {
		return entity.Collection{}, fmt.Errorf("read collection %s: %w", id, err)
	}

	var c entity.Collection
	if err := json.Unmarshal(data, &c); err != nil {
		return entity.Collection{}, fmt.Errorf("parse collection %s: %w", id, err)
	}
	return c, nil
}

// outbound/file.go (lanjutan)

func DeleteCollection(id string) error {
	dir, err := baseDir()
	if err != nil {
		return err
	}

	if err := os.Remove(filepath.Join(dir, id+".json")); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove file: %w", err)
	}

	return RemoveFromIndex(id)
}

// ── Load all ─────────────────────────────────────────────────

// LoadAllCollections membaca semua collection sesuai urutan index.json
func LoadAllCollections() ([]entity.Collection, error) {
	dir, err := baseDir()
	if err != nil {
		return nil, fmt.Errorf("base dir: %w", err)
	}

	index, err := LoadIndex()
	if err != nil {
		return nil, fmt.Errorf("load index: %w", err)
	}

	// index kosong — fallback scan folder
	if len(index.Order) == 0 {
		return scanCollections(dir)
	}

	collections := make([]entity.Collection, 0, len(index.Order))
	for _, id := range index.Order {
		c, err := LoadCollection(id)
		if err != nil {
			// skip file corrupt/hilang, jangan crash semua
			continue
		}
		collections = append(collections, c)
	}

	return collections, nil
}

// scanCollections fallback — baca semua *.json di folder
func scanCollections(dir string) ([]entity.Collection, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir: %w", err)
	}

	var collections []entity.Collection
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "index.json" {
			continue
		}
		if filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		id := strings.TrimSuffix(entry.Name(), ".json")
		c, err := LoadCollection(id)
		if err != nil {
			continue
		}
		collections = append(collections, c)
	}

	return collections, nil
}
