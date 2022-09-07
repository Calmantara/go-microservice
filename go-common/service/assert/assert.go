//go:generate mockgen -source assert.go -destination mock/assert_mock.go -package mock

package serviceassert

import "github.com/google/uuid"

type Assert interface {
	IsZero(value int) bool
	IsEmpty(value string) bool
	IsUUIDEmpty(value uuid.UUID) bool
}

type AssertStruct struct{}

func NewAssert() Assert {
	return &AssertStruct{}
}

func (a *AssertStruct) IsZero(value int) bool {
	return value == 0
}

func (a *AssertStruct) IsEmpty(value string) bool {
	return value == ""
}

func (a *AssertStruct) IsUUIDEmpty(value uuid.UUID) bool {
	return value.String() == "00000000-0000-0000-0000-000000000000"
}
