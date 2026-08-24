package report

// NextVersion returns the version number a revision should use.
func NextVersion(current Report) int {
	return current.Version + 1
}

// IsDraft reports whether a report can still be edited.
func IsDraft(current Report) bool {
	return current.Status == StatusDraft
}
