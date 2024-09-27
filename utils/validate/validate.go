package validate

import (
	"fmt"
	"go-service-template/utils"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Get returns a new validator instance
func Get() *validator.Validate {
	if validate != nil {
		return validate
	}

	validate = validator.New()

	validate.RegisterValidation("gw_uuid", validateGwUuid, false)

	return validate
}

func validateGwUuid(fl validator.FieldLevel) bool {
	gwUuid := fl.Field().String()

	// Remove "gw" prefix from gw_uuid if it exists
	gwUuid = utils.CutGw(gwUuid)

	return utils.IsValidUUID(gwUuid)
}

type ErrorsMap map[string]string

func (e ErrorsMap) Error() string {
	message := ""
	for _, v := range e {
		if message != "" {
			message += "; "
		}
		message += v
	}
	return message
}

// Custom error message mapping
func ErrMessages(model interface{}, err error) ErrorsMap {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorsMap{"error": err.Error()}
	}

	messages := ErrorsMap{}

	for _, err := range errs {
		fieldName := getJSONTag(model, err.Field())

		switch err.Tag() {
		case "required":
			messages[fieldName] = fmt.Sprintf("%s is required", fieldName)
		default:
			messages[fieldName] = fmt.Sprintf("%s is not valid", fieldName)
		}
	}
	return messages
}

// getJSONTag retrieves the JSON tag from the struct field
func getJSONTag(model interface{}, fieldName string) string {
	model = reflect.Indirect(reflect.ValueOf(model)).Interface()
	t := reflect.TypeOf(model)
	field, found := t.FieldByName(fieldName)
	if !found {
		return fieldName // fallback to the original field name if JSON tag is not found
	}
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return fieldName // fallback to the original field name if JSON tag is empty
	}
	return jsonTag
}

// Struct validates a struct
func Struct(s interface{}) error {
	if err := Get().Struct(s); err != nil {
		return ErrMessages(s, err)
	}
	return nil
}
