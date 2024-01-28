package gocast

import "context"

type (
	structWalkerFunc     func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) error
	structWalkerNameFunc func(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) string
)

// StructWalkOptions is the options for StructWalk
type StructWalkOptions struct {
	pathTag       string
	pathExtractor structWalkerNameFunc
}

// PathName returns the path name of the current field
func (w *StructWalkOptions) PathName(ctx context.Context, curObj StructWalkObject, field StructWalkField, path []string) string {
	if w != nil && w.pathExtractor != nil {
		return w.pathExtractor(ctx, curObj, field, path)
	}
	if w != nil && w.pathTag != "" {
		return field.Tag(w.pathTag)
	}
	return field.Name()
}

// WalkOption is the option for StructWalk
type WalkOption func(w *StructWalkOptions)

// WalkWithPathExtractor sets the path extractor for StructWalk
func WalkWithPathExtractor(fn structWalkerNameFunc) WalkOption {
	return func(w *StructWalkOptions) {
		w.pathExtractor = fn
	}
}

// WalkWithPathTag sets the path tag for StructWalk
func WalkWithPathTag(tagName string) WalkOption {
	return func(w *StructWalkOptions) {
		w.pathTag = tagName
	}
}
