package transactions

import "time"

type Transaction struct {
	ID        int        `bson:"id"`
	CreatedOn time.Time  `bson:"createdOn"`
	UpdatedOn *time.Time `bson:"updatedOn"`
	Origin    string     `bson:"origin"`
	SourceID  int        `bson:"sourceAccountId"`
	TargetID  int        `bson:"targetAccountId"`
	Amount    float64    `bson:"amount"`
	Comments  *string    `bson:"comments"`
}
