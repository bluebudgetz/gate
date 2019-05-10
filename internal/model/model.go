package model

import "time"

type Account struct {
	ID        int        `json:"id"`
	CreatedOn time.Time  `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	DeletedOn *time.Time `json:"deletedOn"`
	Name      string     `json:"name"`
	Incoming  float64    `json:"incoming"`
	Outgoing  float64    `json:"outgoing"`
	ParentID  *int       `json:"parentId"`
	Leaf      bool       `json:"leaf"`
}

type Transaction struct {
	ID        int        `json:"id"`
	CreatedOn time.Time  `json:"createdOn"`
	UpdatedOn *time.Time `json:"updatedOn"`
	DeletedOn *time.Time `json:"deletedOn"`
	Origin    string     `json:"origin"`
	SourceID  int        `json:"sourceAccountId"`
	TargetID  int        `json:"targetAccountId"`
	Amount    float64    `json:"amount"`
	Comments  *string    `json:"comments"`
}
