package create

func Create(path, fileText string, showChanges bool) error {
	return (&Creator{}).Create(path, fileText, showChanges)
}
