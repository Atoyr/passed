package orm

import "time"

type Model struct {
	ID       string
	CreateAt time.Time
	UpdateAt time.Time
}
