package create

import (
	"fmt"
	"io"

	"github.com/RRethy/eddie/internal/display"
	"github.com/RRethy/eddie/internal/fileops"
)

type Creator struct {
	fileOps *fileops.FileOps
	display *display.Display
}

func NewCreator(w io.Writer) *Creator {
	return &Creator{
		fileOps: &fileops.FileOps{},
		display: display.New(w),
	}
}

func (c *Creator) Create(path, fileText string, showChanges, showResult bool) error {
	err := c.fileOps.CreateFile(path, fileText)
	if err != nil {
		return err
	}

	if showChanges {
		c.display.ShowNewFileContent(path, fileText)
	}

	if showResult {
		c.display.ShowResult(path, fileText)
	}

	fmt.Printf("Created file: %s (%d bytes)\n", path, len(fileText))
	return nil
}
