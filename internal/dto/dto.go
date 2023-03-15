package dto

type AccountFindFields struct {
	FirstName string `json:"firstName" default:""`
	LastName  string `json:"lastName" default:""`
	Email     string `json:"email" default:""`
	From      int    `json:"from" default:"0"`
	Size      int    `json:"size" default:"10"`
}

type AnimalFindFields struct {
	StartDateTime      string `json:"startDateTime"`
	EndDateTime        string `json:"endDateTime"`
	ChipperId          int    `json:"chipperId"`
	ChippingLocationId int    `json:"chippingLocationId"`
	LifeStatus         string `json:"lifeStatus"`
	Gender             string `json:"gender"`
	From               int    `json:"from" default:"0"`
	Size               int    `json:"size" default:"10"`
}

type AnimalEdit struct {
	OldTypeId int `json:"oldTypeId"`
	NewTypeId int `json:"newTypeId"`
}

type VisitedLocationsFindFields struct {
	StartDateTime string `json:"startDateTime"`
	EndDateTime   string `json:"endDateTime"`
	From          int    `json:"from" default:"0"`
	Size          int    `json:"size" default:"10"`
}
