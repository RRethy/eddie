package undo_edit

func UndoEdit(path string, showChanges, showResult bool) error {
	return (&UndoEditor{}).UndoEdit(path, showChanges, showResult)
}
