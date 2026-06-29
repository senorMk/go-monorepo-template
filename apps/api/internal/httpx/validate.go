package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

// DecodeAndValidate reads a JSON request body into dst, rejecting unknown
// fields (mirrors NestJS forbidNonWhitelisted), then runs struct validation
// (mirrors class-validator). It returns an *AppError suitable for WriteError.
func DecodeAndValidate(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var ute *json.UnmarshalTypeError
		switch {
		case errors.Is(err, io.EOF):
			return BadRequest("Request body is required")
		case errors.As(err, &ute):
			return BadRequest(fmt.Sprintf("Invalid type for field %q", ute.Field))
		default:
			return BadRequest("Invalid request body")
		}
	}

	if err := validate.Struct(dst); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			return BadRequest(formatValidationErrors(verrs))
		}
		return BadRequest("Validation failed")
	}
	return nil
}

// formatValidationErrors joins field errors with "; " to match the NestJS
// ValidationPipe message format.
func formatValidationErrors(verrs validator.ValidationErrors) string {
	msgs := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		field := lowerFirst(fe.Field())
		switch fe.Tag() {
		case "required":
			msgs = append(msgs, fmt.Sprintf("%s is required", field))
		case "email":
			msgs = append(msgs, fmt.Sprintf("%s must be a valid email", field))
		case "min":
			msgs = append(msgs, fmt.Sprintf("%s must be at least %s characters", field, fe.Param()))
		default:
			msgs = append(msgs, fmt.Sprintf("%s is invalid", field))
		}
	}
	return strings.Join(msgs, "; ")
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
