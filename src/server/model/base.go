package model

import "time"

type Model struct {
	Id        int        `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

