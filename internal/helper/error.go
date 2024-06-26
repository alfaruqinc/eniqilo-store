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
		return fmt.Sprintf("minimimum %s is %s characters", field, param)
	case "lte":
		return fmt.Sprintf("maximum %s is %s characters", field, param)
	case "min":
		return fmt.Sprintf("minimum %s is %s", field, param)
	case "max":
		return fmt.Sprintf("maximum %s is %s", field, param)
	case "e164", "url", "validurl", "phonenumber":
		return fmt.Sprintf("wrong %s format", field)
	case "boolean":
		return fmt.Sprintf("%s should be true or false", field)
	case "oneof":
		return fmt.Sprintf("%s should be one of this value %s", field, param)
	case "uuid4":
		return fmt.Sprintf("%s should in uuidv4 format", field)
	case "number":
		return fmt.Sprintf("%s must be number", field)
	}

	return "unhandled validation"
}

func ValidateRequest(err error) domain.MessageErr {
	if err, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range err {
			return domain.NewBadRequestError(msgForTag(fe))
		}
	}
	return domain.NewBadRequestError(err.Error())
}
