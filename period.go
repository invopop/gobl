package gobl

// Period represents two dates with a start and finish.
type Period struct {
	Start *Date `json:"start"`
	End   *Date `json:"end"`
}
