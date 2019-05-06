package model

import "time"

type Account struct {
	ID        int        `json:"id"`
	CreatedOn time.Time  `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	DeletedOn *time.Time `json:"deletedOn"`
	Name      string     `json:"name"`
	Balance   float64    `json:"balance"`
	ParentID  *int       `json:"parentId"`
}

type Transaction struct {
	ID        int       `json:"id"`
	CreatedOn time.Time `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	DeletedOn *time.Time `json:"deletedOn"`
	Origin    string    `json:"origin"`
	SourceID  int       `json:"sourceAccountId"`
	TargetID  int       `json:"targetAccountId"`
	Amount    float64   `json:"amount"`
	Comments  *string   `json:"comments"`
}
