package id

import "github.com/oklog/ulid/v2"

type ULIDGenerator struct{}

func (u *ULIDGenerator) New() string {
	return ulid.Make().String()
}
