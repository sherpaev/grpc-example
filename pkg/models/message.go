package models

import "time"

type Message struct {
	Content   string
	Timestamp time.Time
}
