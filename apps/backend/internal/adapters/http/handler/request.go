package handler

import (
	service "cidadon/internal/application/contract"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func bindRequest(c *gin.Context, request any) bool {
	if err := c.ShouldBind(request); err != nil {
		fields := make([]string, 0)
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				fields = append(fields, jsonFieldName(fieldError.Field()))
			}
		}
		c.Error(service.InvalidInput("invalid request body").WithDetails(map[string]any{"fields": fields}))
		return false
	}
	return true
}

func jsonFieldName(field string) string {
	field = strings.NewReplacer("URL", "Url", "ID", "Id").Replace(field)
	var result strings.Builder
	for index, character := range field {
		if index > 0 && unicode.IsUpper(character) {
			result.WriteByte('_')
		}
		result.WriteRune(unicode.ToLower(character))
	}
	return result.String()
}
