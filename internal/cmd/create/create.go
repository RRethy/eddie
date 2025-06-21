package create

func Create(path, fileText string, showChanges, showResult bool) error {
	return (&Creator{}).Create(path, fileText, showChanges, showResult)
}
