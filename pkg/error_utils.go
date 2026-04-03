package pkg

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

/*
function ini buat ngecek validasi di c.ShouldBindJson(&req) tuh error validasi atau bukan
liat dari argumen kedua yang di return. Kalo true ->

	artinya dia validasi error, dan bakal balikin : validator.ValidationErrors

kalo false -> berarti salah kirim json body nya, either json rusak or smth
*/
func ParseValidationErrors(err error) (map[string]string, bool) {
	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {
		return nil, false
	}

	out := make(map[string]string)

	for _, fe := range ve {
		field := strings.ToLower(fe.Field())
		out[field] = fe.Error()
	}

	return out, true
}
