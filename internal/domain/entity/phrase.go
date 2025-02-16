package entity

import "time"

type Phrase struct {
	ID        uint      `db:"id"`
	UserID    uint      `db:"user_id"`
	Phrase    string    `db:"phrase"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
