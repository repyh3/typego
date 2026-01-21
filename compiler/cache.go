package compiler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const CacheVersion = "v1"

type CacheEntry struct {
	Hash      string   `json:"hash"`
	Imports   []string `json:"imports"`
	JS        string   `json:"js"`
	SourceMap string   `json:"source_map"`
}

func getCacheDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(wd, ".typego", "cache")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func computeHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func getCachePath(filePath string) (string, error) {
	hash, err := computeHash(filePath)
	if err != nil {
		return "", err
	}

	dir, err := getCacheDir()
	if err != nil {
		return "", err
	}

	// Filename: <filename-base>_<hash>_<version>.json
	base := filepath.Base(filePath)
	return filepath.Join(dir, fmt.Sprintf("%s_%s_%s.json", base, hash, CacheVersion)), nil
}

func CheckCache(entryPoint string) (*Result, error) {
	path, err := getCachePath(entryPoint)
	if err != nil {
		return nil, err // Ignore cache errors, force build
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil // Cache miss
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, nil // Corrupt cache
	}

	return &Result{
		JS:        entry.JS,
		SourceMap: entry.SourceMap,
		Imports:   entry.Imports,
	}, nil
}

func SaveCache(entryPoint string, res *Result) error {
	path, err := getCachePath(entryPoint)
	if err != nil {
		return err
	}

	entry := CacheEntry{
		Imports:   res.Imports,
		JS:        res.JS,
		SourceMap: res.SourceMap,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
