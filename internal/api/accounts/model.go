package accounts

import "time"

type Account struct {
	ID         int        `json:"id"`
	CreatedOn  time.Time  `json:"createdOn"`
	UpdatedOn  *time.Time `json:"updatedOn"`
	Name       string     `json:"name"`
	ChildCount int        `json:"childCount"`
	Incoming   float64    `json:"incoming"`
	Outgoing   float64    `json:"outgoing"`
	Balance    float64    `json:"balance"`
}
