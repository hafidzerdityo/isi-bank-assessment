package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)


var (
	CURRENT_DATE = time.Now().Format("2006-01-02")
)

func ValidateStruct(s interface{}) (errMsg string, err error) {
    validate := validator.New()
    if err = validate.Struct(s); err != nil {
        errs := make(map[string]string)
        for _, err := range err.(validator.ValidationErrors) {
            errs[err.Field()] = err.Tag()
        }
		errMsg = "Validation failed:"
        for field, tag := range errs {
            errMsg += fmt.Sprintf(" %s (%s),", field, tag)
        }
        errMsg = errMsg[:len(errMsg)-1]


        return errMsg, err
    }
    return "", nil
}

func IsDigit(str string) bool {
	for _, char := range str {
		if !unicode.IsDigit(char) {
			return false
		}
	}
	return true
}

func ValidatePhoneNumber(phoneNumber string) (string, bool) {
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")
	if strings.HasPrefix(phoneNumber, "+62") || strings.HasPrefix(phoneNumber, "62") {
		// Replace "+62" or "62" with "0"
		phoneNumber = strings.Replace(phoneNumber, "+62", "0", 1)
		phoneNumber = strings.Replace(phoneNumber, "62", "0", 1)
	}

	// This regex allows for 10 to 12 digits after the leading zero
	match, _ := regexp.MatchString(`^08\d{9,11}$`, phoneNumber)
	return phoneNumber, match
}

func GenerateShortUUID() string {
	// Generate a UUID
	id := uuid.New()

	// Remove hyphens from the UUID
	idString := strings.ReplaceAll(id.String(), "-", "")

	// Extract first 10 digits
	shortUUID := idString[:10]

	return shortUUID
}