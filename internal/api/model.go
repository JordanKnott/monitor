package api

type UserAccount struct {
	ID       string `json:"id"`
	Nicename string `json:"nicename"`
	Photo    string `json:"photo"`
	Email    string `json:"email"`
}
