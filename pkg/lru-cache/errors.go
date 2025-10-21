package lrucache

import "errors"

var (
	ErrNotFound       = errors.New("item not found")
	ErrExpired        = errors.New("item expired")
	ErrUnexpectedType = errors.New("unexpected type")
)
