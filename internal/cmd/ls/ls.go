package ls

func Ls(path string) error {
	return (&Lister{}).Ls(path)
}
