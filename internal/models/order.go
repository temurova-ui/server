package models

import "time"

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ItemName  string    `json:"item_name"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Status struct{
	New string `json:"new"`
	Paid string `json:"paid"`
}