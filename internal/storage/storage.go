package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Inventory struct {
	Timestamp time.Time       `json:"timestamp"`
	Items     []InventoryItem `json:"items"`
}

func (i *Inventory) String() string {
	return fmt.Sprintf("%+v", *i)
}

type InventoryItem struct {
	IconURL           string `json:"icon_url"`
	ActionInspectLink string `json:"inspect_url"`
	Name              string `json:"name"`
	NameColor         string `json:"name_color"`
	MarketName        string `json:"market_name"`
	MarketHashName    string `json:"market_hash_name"`
	MarketInspectLink string `json:"market_inspect_link"`
	Marketable        bool   `json:"marketable"`
	Tradable          bool   `json:"tradable"`
	Amount            int    `json:"amount"`
	Price             int    `json:"price"`
	Currency          string `json:"currency"`
}

func (i InventoryItem) String() string {
	return fmt.Sprintf(
		"InventoryItem{IconURL: %s, ActionInspectLink: %s, Name: %s, NameColor: %s, MarketName: %s, MarketHashName: %s, MarketInspectLink: %s, Marketable: %t, Tradable: %t, Amount: %d, Price: %d, Currency: %s}",
		i.IconURL,
		i.ActionInspectLink,
		i.Name,
		i.NameColor,
		i.MarketName,
		i.MarketHashName,
		i.MarketInspectLink,
		i.Marketable,
		i.Tradable,
		i.Amount,
		i.Price,
		i.Currency,
	)
}

func Write(projectName string, items []InventoryItem) error {
	if projectName == "" {
		return errors.New("project name cannot be empty")
	}

	if len(items) == 0 {
		return errors.New("no items to write")
	}

	i := &Inventory{
		Timestamp: time.Now(),
		Items:     items,
	}

	storageFile, err := createStorageFile(projectName)
	if err != nil {
		return fmt.Errorf("failed to create storage file: %w", err)
	}
	defer storageFile.Close()

	err = json.NewEncoder(storageFile).Encode(i)
	if err != nil {
		return fmt.Errorf("failed to encode inventory to storage file: %w", err)
	}

	return nil
}

var storageFileTimestamp = time.Now().Format("2006-01-02_15-04-05")

func createStorageFile(projectName string) (*os.File, error) {
	storageDir, err := setupStorageDir(projectName)
	if err != nil {
		return nil, fmt.Errorf("failed to setup storage directory: %w", err)
	}

	storageFilePath := filepath.Join(storageDir, "storage_"+storageFileTimestamp+".json")

	file, err := os.Create(storageFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage file %s: %w", storageFilePath, err)
	}

	return file, nil
}

func setupStorageDir(projectName string) (string, error) {
	storagesDir, err := setupStoragesDir()
	if err != nil {
		return "", fmt.Errorf("failed to setup storages directory: %w", err)
	}

	storageDir := filepath.Join(storagesDir, projectName)

	err = os.MkdirAll(storageDir, 0700)
	if err != nil {
		return "", fmt.Errorf("failed to create storage directory %s: %w", storageDir, err)
	}

	return storageDir, nil
}

func setupStoragesDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	storagesDir := filepath.Join(home, ".dropawp", "storages")

	err = os.MkdirAll(storagesDir, 0700)
	if err != nil {
		return "", fmt.Errorf("failed to create storages directory %s: %w", storagesDir, err)
	}

	return storagesDir, nil
}
