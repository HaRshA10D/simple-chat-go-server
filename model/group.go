package model

import (
	"time"
)

type Group struct {
	ID string
	Name string
	LastActivityAt time.Time
}
