/*
Package jsoncache offers easy-to-use file based caching through json to easily cache data with minimal effort.
*/
package jsoncache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Cache is a struct that manages file-based caching
type Cache struct {
	directory string
	timeout   time.Duration
}

// NewCache creates a new Cache instance
func NewCache(directory string, timeout time.Duration) *Cache {
	return &Cache{directory: directory, timeout: timeout}
}

// CacheData is wraps the data with a timestamp for expiry
type CacheData struct {
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

// Save saves a JSON object to a file identified by a unique ID
func (c *Cache) Save(id string, data interface{}) error {
	if err := os.MkdirAll(c.directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}

	cacheData := CacheData{
		Timestamp: time.Now(),
		Data:      data,
	}

	filePath := filepath.Join(c.directory, id+".json")

	jsonData, err := json.MarshalIndent(cacheData, "", "   ")
	if err != nil {
		return fmt.Errorf("failed to serialize data: %v", err)
	}

	if err := os.WriteFile(filePath, jsonData, 0o644); err != nil {
		return fmt.Errorf("failed to write data to file: %v", err)
	}

	return nil
}

// Load loads a JSON object from a file identified by a unique ID
func (c *Cache) Load(id string, data interface{}) error {
	filePath := filepath.Join(c.directory, id+".json")

	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read data from file: %v", err)
	}

	var cacheData CacheData

	if err = json.Unmarshal(jsonData, &cacheData); err != nil {
		return fmt.Errorf("failed to deserialize data: %v", err)
	}

	// Check if the data is within the timeout
	if time.Since(cacheData.Timestamp) > c.timeout {
		return fmt.Errorf("cached data has expired")
	}

	dataBytes, err := json.Marshal(cacheData.Data)
	if err != nil {
		return fmt.Errorf("failed to re-marshal inner data: %v", err)
	}

	if err := json.Unmarshal(dataBytes, data); err != nil {
		return fmt.Errorf("failed to unmarshal inner data: %v", err)
	}

	return nil
}
