package str_replace

func StrReplace(path, oldStr, newStr string) error {
	return (&Replacer{}).StrReplace(path, oldStr, newStr)
}
