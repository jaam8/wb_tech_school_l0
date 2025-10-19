package lru_cache

import "errors"

var (
	ErrNotFound = errors.New("item not found")
	ErrExpired  = errors.New("item expired")
)
