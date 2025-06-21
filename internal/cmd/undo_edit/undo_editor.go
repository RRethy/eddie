package undo_edit

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type UndoEditor struct{}

type EditRecord struct {
	EditType    string    `json:"edit_type"`     // "str_replace" or "insert"
	OldContent  string    `json:"old_content"`   // For str_replace: old string, for insert: ""
	NewContent  string    `json:"new_content"`   // For str_replace: new string, for insert: inserted line
	Position    int       `json:"position"`      // For insert: line number, for str_replace: -1
	Timestamp   time.Time `json:"timestamp"`     // When edit was made
	FileModTime time.Time `json:"file_mod_time"` // File modification time after edit
}

type EditHistory struct {
	FilePath string       `json:"file_path"`
	Edits    []EditRecord `json:"edits"`
}

func (u *UndoEditor) UndoEdit(path string) error {
	// Check if original file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	// Get edit file path
	editPath, err := u.getEditFilePath(path)
	if err != nil {
		return fmt.Errorf("get edit file path: %w", err)
	}

	// Read edit history
	editHistory, err := u.readEditHistory(editPath)
	if err != nil {
		return fmt.Errorf("read edit history %s: %w", editPath, err)
	}

	if len(editHistory.Edits) == 0 {
		return fmt.Errorf("no edit records found for %s", path)
	}

	// Get the most recent edit (last in array)
	lastEditIndex := len(editHistory.Edits) - 1
	editRecord := editHistory.Edits[lastEditIndex]

	// Validate file modification time
	if !info.ModTime().Equal(editRecord.FileModTime) {
		return fmt.Errorf("file has been modified since last tracked edit (expected: %v, actual: %v)", 
			editRecord.FileModTime.Format(time.RFC3339), info.ModTime().Format(time.RFC3339))
	}

	// Apply reverse edit
	err = u.applyReverseEdit(path, &editRecord)
	if err != nil {
		return fmt.Errorf("apply reverse edit: %w", err)
	}

	// Remove the last edit from history
	editHistory.Edits = editHistory.Edits[:lastEditIndex]

	// If there are remaining edits, update the mod time of the most recent one
	if len(editHistory.Edits) > 0 {
		// Get new file mod time after undo
		newInfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat file after undo: %w", err)
		}
		
		// Update the mod time of the now-most-recent edit
		editHistory.Edits[len(editHistory.Edits)-1].FileModTime = newInfo.ModTime()
	}

	// Save updated history or remove file if empty
	if len(editHistory.Edits) == 0 {
		err = os.Remove(editPath)
		if err != nil {
			return fmt.Errorf("remove empty edit file %s: %w", editPath, err)
		}
	} else {
		err = u.writeEditHistory(editPath, editHistory)
		if err != nil {
			return fmt.Errorf("write updated edit history: %w", err)
		}
	}

	fmt.Printf("Undid %s edit in %s\n", editRecord.EditType, path)
	return nil
}


// RecordEdit records an edit operation for potential undo
func (u *UndoEditor) RecordEdit(path, editType, oldContent, newContent string, position int) error {
	// Get edit file path
	editPath, err := u.getEditFilePath(path)
	if err != nil {
		return fmt.Errorf("get edit file path: %w", err)
	}

	// Ensure edit directory exists
	editDir := filepath.Dir(editPath)
	err = os.MkdirAll(editDir, 0755)
	if err != nil {
		return fmt.Errorf("create edit directory %s: %w", editDir, err)
	}

	// Get file modification time after the edit
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file %s: %w", path, err)
	}

	// Create new edit record
	newEdit := EditRecord{
		EditType:    editType,
		OldContent:  oldContent,
		NewContent:  newContent,
		Position:    position,
		Timestamp:   time.Now(),
		FileModTime: info.ModTime(),
	}

	// Read existing edit history or create new one
	editHistory, err := u.readEditHistory(editPath)
	if err != nil {
		// File doesn't exist, create new history
		editHistory = &EditHistory{
			FilePath: path,
			Edits:    []EditRecord{},
		}
	}

	// Append new edit to history
	editHistory.Edits = append(editHistory.Edits, newEdit)

	// Write updated history
	err = u.writeEditHistory(editPath, editHistory)
	if err != nil {
		return fmt.Errorf("write edit history: %w", err)
	}

	return nil
}

// getEditDir returns the edit directory path using XDG cache directory
func (u *UndoEditor) getEditDir() (string, error) {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get user home directory: %w", err)
		}
		cacheDir = filepath.Join(homeDir, ".cache")
	}

	return filepath.Join(cacheDir, "eddie", "edits"), nil
}

// getEditFilePath creates the edit file path based on the original file name
func (u *UndoEditor) getEditFilePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("get absolute path: %w", err)
	}

	editDir, err := u.getEditDir()
	if err != nil {
		return "", fmt.Errorf("get edit directory: %w", err)
	}

	// Create safe filename from absolute path
	safeName := u.createSafeFilename(absPath)
	editFileName := safeName + ".json"
	
	return filepath.Join(editDir, editFileName), nil
}

// createSafeFilename converts a file path to a safe filename for storage
func (u *UndoEditor) createSafeFilename(path string) string {
	// Replace path separators and other problematic characters
	safe := strings.ReplaceAll(path, string(filepath.Separator), "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	safe = strings.ReplaceAll(safe, " ", "_")
	safe = strings.ReplaceAll(safe, ".", "_")
	
	// If the filename would be too long, use a hash
	if len(safe) > 200 {
		hash := sha256.Sum256([]byte(path))
		return fmt.Sprintf("file_%x", hash)[:50]
	}
	
	return safe
}

// readEditHistory reads edit history from file
func (u *UndoEditor) readEditHistory(editPath string) (*EditHistory, error) {
	data, err := os.ReadFile(editPath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var history EditHistory
	err = json.Unmarshal(data, &history)
	if err != nil {
		return nil, fmt.Errorf("unmarshal JSON: %w", err)
	}

	return &history, nil
}

// writeEditHistory writes edit history to file
func (u *UndoEditor) writeEditHistory(editPath string, history *EditHistory) error {
	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal JSON: %w", err)
	}

	err = os.WriteFile(editPath, data, 0644)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
