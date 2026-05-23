package dto

type QueryEntry struct {
	Query string `json:"query"`
	Count int    `json:"count"`
}

type TopQueriesResponse struct {
	Top []QueryEntry `json:"top"`
}
