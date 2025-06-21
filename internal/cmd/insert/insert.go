package insert

func Insert(path, insertLine, newStr string) error {
	return (&Inserter{}).Insert(path, insertLine, newStr)
}
