package helper

import (
	"fmt"
	"strings"
)

func MsgForTag(field, tag, param string) string {
	field = strings.ToLower(field)

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "gte":
		return fmt.Sprintf("%s should greater than or equal to %s", field, param)
	case "lte":
		return fmt.Sprintf("%s should less than or equal to %s", field, param)
	}

	return "unhandled validation"
}
