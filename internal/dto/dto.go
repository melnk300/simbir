package dto

type FindFields struct {
	FirstName string `json:"firstName" default:""`
	LastName  string `json:"lastName" default:""`
	Email     string `json:"email" default:""`
	From      int    `json:"from" default:"0"`
	Size      int    `json:"size" default:"10"`
}
