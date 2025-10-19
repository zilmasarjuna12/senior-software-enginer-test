package model

type FetchEventsRequest struct {
	Severity int      `json:"severity"`
	Tags     []string `json:"tags"`
}
