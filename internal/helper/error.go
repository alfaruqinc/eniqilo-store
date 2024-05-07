package helper

import (
	"eniqilo-store/internal/domain"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func msgForTag(fe validator.FieldError) string {
	field := fe.Field()
	param := fe.Param()

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "gte":
		return fmt.Sprintf("%s should greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s should less than or equal to %s", field, param)
	case "e164":
		return fmt.Sprintf("wrong %s format", field)
	}

	return "unhandled validation"
}

func ValidateRequest(err error) domain.MessageErr {
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range err {
			return domain.NewBadRequestError(msgForTag(fe))
		}
	}
	return domain.NewBadRequestError("unhandled validation")
}
