package id

import "github.com/oklog/ulid/v2"

// ULIDGenerator generates ULIDs for unique identifiers.
type ULIDGenerator struct{}

func (u *ULIDGenerator) New() string {
	return ulid.Make().String()
}
