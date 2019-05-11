package model

import "time"

type Account struct {
	ID         int        `json:"id"`
	CreatedOn  time.Time  `json:"createdOn"`
	UpdatedOn  *time.Time `json:"updatedOn"`
	Name       string     `json:"name"`
	ChildCount int        `json:"childCount"`
	Incoming   float64    `json:"incoming"`
	Outgoing   float64    `json:"outgoing"`
}

type Transaction struct {
	ID        int        `json:"id"`
	CreatedOn time.Time  `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	Origin    string     `json:"origin"`
	SourceID  int        `json:"sourceAccountId"`
	TargetID  int        `json:"targetAccountId"`
	Amount    float64    `json:"amount"`
	Comments  *string    `json:"comments"`
}
