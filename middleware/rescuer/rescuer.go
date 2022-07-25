package rescuer

import (
	"errors"
	"net/http"

	"homework/utility/mytype"

	"github.com/uptrace/bunrouter"
)

func NewMiddleware() bunrouter.MiddlewareFunc {
	return middlewareFunc
}

func middlewareFunc(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) (err error) {
		defer func() {
			if r := recover(); r != nil {
				switch x := r.(type) {
				case string:
					err = errors.New(x)
				case error:
					err = x
				default:
					err = errors.New("Unknown panic")
				}
			}

			if err != nil {
				if apiError, ok := err.(mytype.ApiError); ok {
					http.Error(w, apiError.Error(), apiError.Status)
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}
		}()

		return next(w, req)
	}
}
