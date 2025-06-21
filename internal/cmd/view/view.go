package view

func View(path, viewRange string) error {
	return (&Viewer{}).View(path, viewRange)
}
