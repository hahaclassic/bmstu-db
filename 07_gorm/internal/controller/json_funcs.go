package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"github.com/hahaclassic/databases/07_gorm/internal/models"
)

func (c *Controller) ExportUsersToJSON(ctx context.Context) {
	var fileName string

	fmt.Print("Enter the file name for export (default: ./data/users.json): ")
	_, err := fmt.Scanln(&fileName)
	if err != nil || fileName == "" {
		fileName = "./data/users.json"
	}

	jsonData, err := c.storage.ExportUsersToJSON(ctx)
	if err != nil {
		slog.Error("Error exporting users to JSON", "err", err)
		return
	}

	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		slog.Error("Error writing JSON to file", "err", err)
		return
	}

	fmt.Printf("Users successfully exported to '%s'\n", fileName)
}

func (c *Controller) ImportUsersFromJSON(ctx context.Context) {
	var fileName string

	fmt.Print("Enter the file name for import (default: ./data/users.json): ")
	_, err := fmt.Scanln(&fileName)
	if err != nil || fileName == "" {
		fileName = "./data/users.json"
	}

	jsonData, err := os.ReadFile(fileName)
	if err != nil {
		slog.Error("Error reading JSON from file", "err", err)
		return
	}

	var users []*models.User
	if err := json.Unmarshal(jsonData, &users); err != nil {
		slog.Error("failed to parse JSON", "err", err)
		return
	}

	err = c.storage.ImportUsers(ctx, users)
	if err != nil {
		slog.Error("Error importing users from JSON", "err", err)
		return
	}

	fmt.Printf("Successfully imported.\n")
}
