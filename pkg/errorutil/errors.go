package errorutil

import (
	"fmt"

	"github.com/pkg/errors"
)

func InternalError(v interface{}) error {
	switch v := v.(type) {
	case string:
		return errors.New(v)
	case error:
		return errors.New(v.Error())
	default:
		return errors.New(fmt.Sprintf("unknown error value: %+v", v))
	}
}
