package undo_edit

import "os"

func UndoEdit(path string, showChanges, showResult bool, count int) error {
	return NewUndoEditor(os.Stdout).UndoEdit(path, showChanges, showResult, count)
}
