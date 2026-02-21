package id

// Generator defines an interface for generating unique IDs.
type Generator interface {
	New() string
}
