package domain

import "time"

type User struct {
	ID        string
	Email     string
	Name      string
	Roles     []string
	CreatedAt time.Time
}
