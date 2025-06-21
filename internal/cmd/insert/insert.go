package insert

import "os"

func Insert(path, insertLine, newStr string, showChanges, showResult bool) error {
	return NewInserter(os.Stdout).Insert(path, insertLine, newStr, showChanges, showResult)
}
