package batch

type BatchRequest struct {
	Operations []Operation `json:"operations"`
}

type Operation struct {
	Type string `json:"type"`

	Path      string `json:"path,omitempty"`
	ViewRange string `json:"view_range,omitempty"`

	OldStr      string `json:"old_str,omitempty"`
	NewStr      string `json:"new_str,omitempty"`
	ShowChanges bool   `json:"show_changes,omitempty"`
	ShowResult  bool   `json:"show_result,omitempty"`

	Content string `json:"content,omitempty"`

	InsertLine int    `json:"insert_line,omitempty"`
	Count      int    `json:"count,omitempty"`
	TreeQuery  string `json:"tree_sitter_query,omitempty"`
	Pattern    string `json:"pattern,omitempty"`
}

type BatchResponse struct {
	Results []OperationResult `json:"results"`
}

type OperationResult struct {
	Operation Operation `json:"operation"`
	Success   bool      `json:"success"`
	Output    string    `json:"output"`
	Error     *string   `json:"error"`
}
