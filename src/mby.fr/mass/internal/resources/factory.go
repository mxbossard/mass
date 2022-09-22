package resources

import (
	"fmt"
	"reflect"
)

/*
type Res interface {
	Env | Project | Image
}

type ResPtr interface {
	*Env | *Project | *Image
}
*/

func FromPath[T Resource] (path string) (res T, err error) {
	r, err := Read(path)
	if err != nil {
		return
	}

	res, ok := r.(T)
	if reflect.ValueOf(res).Kind() != reflect.Ptr {
		// Expect resource value
		// In this case, r is a pointer and we want to return a value, but the type cast don't return ok.
		resPtrType := reflect.PointerTo(reflect.TypeOf(res))
		if reflect.TypeOf(r) == resPtrType {
			// Right type so res was rightly type cast
			return res, err
		}
	}

	if !ok {
		err = fmt.Errorf("bad resource type for path %s. Expected type %T but got %T", path, res, r)
	}
	return res, err
}
