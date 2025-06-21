package str_replace

func StrReplace(path, oldStr, newStr string, showChanges bool) error {
	return (&Replacer{}).StrReplace(path, oldStr, newStr, showChanges)
}
