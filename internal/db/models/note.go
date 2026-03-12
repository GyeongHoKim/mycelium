// Package models define note model that db needs
package models

import "time"

type Note struct {
	Path        string
	ContentHash uint64
	VectorID    string
	UpdatedAt   time.Time
}
