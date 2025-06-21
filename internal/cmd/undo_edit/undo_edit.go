package undo_edit

import "os"

func UndoEdit(path string, showChanges, showResult bool) error {
	return NewUndoEditor(os.Stdout).UndoEdit(path, showChanges, showResult)
}
