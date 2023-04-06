package repository

import "titmouse/lib/log"

type Repository struct {
}

func (customR Repository) log() *log.Log {
	return log.Api()
}
