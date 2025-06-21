package str_replace

import "os"

func StrReplace(path, oldStr, newStr string, showChanges, showResult bool) error {
	return NewReplacer(os.Stdout).StrReplace(path, oldStr, newStr, showChanges, showResult)
}
