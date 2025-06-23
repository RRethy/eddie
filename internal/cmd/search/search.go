package search

func Search(path, query string) error {
	return (&Searcher{}).Search(path, query)
}
