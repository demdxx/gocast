package gocast

import "context"

// CastSetter interface from some type into the specific value
type CastSetter interface {
	CastSet(ctx context.Context, v any) error
}
