package domain

import "errors"

var (
	ErrNotFound      = errors.New("your requested asset info is not found")
	ErrNoEnoughAsset = errors.New("your asset is not enough")
)
