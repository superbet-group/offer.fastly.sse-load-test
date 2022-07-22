package models

// Match is a simplified match model received from the offer server
type Match struct {
	ID int64 `json:"mi"`
}

// Offer is a simplified offer model received from the offer server
type Offer struct {
	Error bool
	Data  []Match `json:"data"`
}
