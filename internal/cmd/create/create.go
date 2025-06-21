package create

func Create(path, fileText string) error {
	return (&Creator{}).Create(path, fileText)
}