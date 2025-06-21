package insert

func Insert(path, insertLine, newStr string, showChanges, showResult bool) error {
	return (&Inserter{}).Insert(path, insertLine, newStr, showChanges, showResult)
}
