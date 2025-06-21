package undo_edit

func UndoEdit(path string) error {
	return (&UndoEditor{}).UndoEdit(path)
}
