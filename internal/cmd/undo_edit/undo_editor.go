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

func (u *UndoEditor) UndoEdit(path string, showChanges, showResult bool) error {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	editPath, err := u.getEditFilePath(path)
	if err != nil {
		return fmt.Errorf("get edit file path: %w", err)
	}

	editHistory, err := u.readEditHistory(editPath)
	if err != nil {
		return fmt.Errorf("read edit history %s: %w", editPath, err)
	}

	if len(editHistory.Edits) == 0 {
		return fmt.Errorf("no edit records found for %s", path)
	}

	lastEditIndex := len(editHistory.Edits) - 1
	editRecord := editHistory.Edits[lastEditIndex]

	if !info.ModTime().Equal(editRecord.FileModTime) {
		return fmt.Errorf("file has been modified since last tracked edit (expected: %v, actual: %v)",
			editRecord.FileModTime.Format(time.RFC3339), info.ModTime().Format(time.RFC3339))
	}

	var beforeContent string
	if showChanges {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file before undo: %w", err)
		}
		beforeContent = string(content)
	}

	err = u.applyReverseEdit(path, &editRecord)
	if err != nil {
		return fmt.Errorf("apply reverse edit: %w", err)
	}

	if showChanges || showResult {
		afterContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file after undo: %w", err)
		}
		if showChanges {
			u.showDiff(path, beforeContent, string(afterContent))
		}
		if showResult {
			u.showResult(path, string(afterContent))
		}
	}

	editHistory.Edits = editHistory.Edits[:lastEditIndex]

	if len(editHistory.Edits) > 0 {
		newInfo, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("stat file after undo: %w", err)
		}

		editHistory.Edits[len(editHistory.Edits)-1].FileModTime = newInfo.ModTime()
	}

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

func (u *UndoEditor) RecordEdit(path, editType, oldContent, newContent string, position int) error {
	editPath, err := u.getEditFilePath(path)
	if err != nil {
		return fmt.Errorf("get edit file path: %w", err)
	}

	editDir := filepath.Dir(editPath)
	err = os.MkdirAll(editDir, 0755)
	if err != nil {
		return fmt.Errorf("create edit directory %s: %w", editDir, err)
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stat file %s: %w", path, err)
	}

	newEdit := EditRecord{
		EditType:    editType,
		OldContent:  oldContent,
		NewContent:  newContent,
		Position:    position,
		Timestamp:   time.Now(),
		FileModTime: info.ModTime(),
	}

	editHistory, err := u.readEditHistory(editPath)
	if err != nil {
		editHistory = &EditHistory{
			FilePath: path,
			Edits:    []EditRecord{},
		}
	}

	editHistory.Edits = append(editHistory.Edits, newEdit)

	err = u.writeEditHistory(editPath, editHistory)
	if err != nil {
		return fmt.Errorf("write edit history: %w", err)
	}

	return nil
}

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

func (u *UndoEditor) getEditFilePath(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("get absolute path: %w", err)
	}

	editDir, err := u.getEditDir()
	if err != nil {
		return "", fmt.Errorf("get edit directory: %w", err)
	}

	safeName := u.createSafeFilename(absPath)
	editFileName := safeName + ".json"

	return filepath.Join(editDir, editFileName), nil
}

func (u *UndoEditor) createSafeFilename(path string) string {
	safe := strings.ReplaceAll(path, string(filepath.Separator), "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	safe = strings.ReplaceAll(safe, " ", "_")
	safe = strings.ReplaceAll(safe, ".", "_")

	if len(safe) > 200 {
		hash := sha256.Sum256([]byte(path))
		return fmt.Sprintf("file_%x", hash)[:50]
	}

	return safe
}

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

func (u *UndoEditor) showDiff(path, before, after string) {
	fmt.Printf("\nChanges in %s:\n", path)
	fmt.Println("--- Before undo")
	fmt.Println("+++ After undo")
	
	beforeLines := strings.Split(before, "\n")
	afterLines := strings.Split(after, "\n")
	
	maxLines := len(beforeLines)
	if len(afterLines) > maxLines {
		maxLines = len(afterLines)
	}
	
	for i := 0; i < maxLines; i++ {
		beforeLine := ""
		afterLine := ""
		
		if i < len(beforeLines) {
			beforeLine = beforeLines[i]
		}
		if i < len(afterLines) {
			afterLine = afterLines[i]
		}
		
		if beforeLine != afterLine {
			if beforeLine != "" {
				fmt.Printf("-%s\n", beforeLine)
			}
			if afterLine != "" {
				fmt.Printf("+%s\n", afterLine)
			}
		}
	}
	fmt.Println()
}

func (u *UndoEditor) showResult(path, content string) {
	fmt.Printf("\nResult of %s:\n", path)
	fmt.Println(content)
}
