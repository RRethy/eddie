package create

import "os"

func Create(path, fileText string, showChanges, showResult bool) error {
	return NewCreator(os.Stdout).Create(path, fileText, showChanges, showResult)
}
