package domain

import "errors"

var (
	ErrNotFoundAsset = errors.New("your requested asset info is not found")
	ErrNoEnoughAsset = errors.New("your asset is not enough")

	ErrInnerNilPointer = errors.New("inner nil pointer")
)
