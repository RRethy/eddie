package insert

func Insert(path, insertLine, newStr string, showChanges bool) error {
	return (&Inserter{}).Insert(path, insertLine, newStr, showChanges)
}
