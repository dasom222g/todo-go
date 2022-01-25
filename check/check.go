package check

import (
	"net/http"

	"github.com/unrolled/render"
)

func IsError(err error, rd *render.Render, w http.ResponseWriter, code int) bool {
	if err != nil {
		rd.Text(w, code, err.Error())
		return true
	}
	return false
}
