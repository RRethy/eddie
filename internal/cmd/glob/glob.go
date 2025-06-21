package glob

func Glob(pattern, path string) error {
	return (&Globber{}).Glob(pattern, path)
}