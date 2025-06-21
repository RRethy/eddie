package undo_edit

func UndoEdit(path string, showChanges bool) error {
	return (&UndoEditor{}).UndoEdit(path, showChanges)
}
