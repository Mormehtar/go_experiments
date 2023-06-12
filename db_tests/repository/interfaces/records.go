package interfaces

import (
	"errors"
	"time"
)

var (
	NameIsUsedAlready  = errors.New("name is used already")
	NameIsNotFound     = errors.New("name is not found")
	PropertyIsNotFound = errors.New("property is not found")
)

type Name struct {
	Id        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Property struct {
	Id        int64
	NameId    int64
	Key       string
	Value     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
