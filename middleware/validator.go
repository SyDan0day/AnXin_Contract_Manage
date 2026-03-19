package middleware

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type ValidationRule struct {
	Field    string
	Rules    []string
	MinLen   int
	MaxLen   int
	Pattern  string
	Required bool
}

var (
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	phoneRegex    = regexp.MustCompile(`^1[3-9]\d{9}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	alphanumeric  = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

func ValidateInput(rules []ValidationRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errors []string

		for _, rule := range rules {
			var value string

			switch c.Request.Method {
			case "GET":
				value = c.Query(rule.Field)
				if value == "" {
					value = c.Param(rule.Field)
				}
			case "POST", "PUT", "DELETE":
				value = c.PostForm(rule.Field)
				if value == "" {
					jsonValue, exists := c.Get(rule.Field)
					if exists {
						value = strconv.FormatBool(jsonValue.(bool))
					}
				}
			}

			if rule.Required && strings.TrimSpace(value) == "" {
				errors = append(errors, rule.Field+" is required")
				continue
			}

			if strings.TrimSpace(value) == "" {
				continue
			}

			for _, r := range rule.Rules {
				switch r {
				case "email":
					if !emailRegex.MatchString(value) {
						errors = append(errors, rule.Field+" must be a valid email")
					}
				case "phone":
					if !phoneRegex.MatchString(value) {
						errors = append(errors, rule.Field+" must be a valid phone number")
					}
				case "username":
					if !usernameRegex.MatchString(value) {
						errors = append(errors, rule.Field+" must be 3-20 characters (letters, numbers, underscore)")
					}
				case "alphanumeric":
					if !alphanumeric.MatchString(value) {
						errors = append(errors, rule.Field+" must contain only letters and numbers")
					}
				case "numeric":
					if _, err := strconv.Atoi(value); err != nil {
						errors = append(errors, rule.Field+" must be a number")
					}
				}
			}

			if rule.MinLen > 0 && len(value) < rule.MinLen {
				errors = append(errors, rule.Field+" must be at least "+strconv.Itoa(rule.MinLen)+" characters")
			}

			if rule.MaxLen > 0 && len(value) > rule.MaxLen {
				errors = append(errors, rule.Field+" must be at most "+strconv.Itoa(rule.MaxLen)+" characters")
			}

			if rule.Pattern != "" {
				pattern := regexp.MustCompile(rule.Pattern)
				if !pattern.MatchString(value) {
					errors = append(errors, rule.Field+" format is invalid")
				}
			}
		}

		if len(errors) > 0 {
			c.JSON(400, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func SanitizeInput(input string) string {
	input = strings.TrimSpace(input)

	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#x27;",
		"/", "&#x2F;",
		"`", "&#96;",
	)

	return replacer.Replace(input)
}

func ValidateAndSanitizeInput(rules []ValidationRule) gin.HandlerFunc {
	return func(c *gin.Context) {
		var errors []string

		for _, rule := range rules {
			var value string

			switch c.Request.Method {
			case "GET":
				value = c.Query(rule.Field)
				if value == "" {
					value = c.Param(rule.Field)
				}
			default:
				value = c.PostForm(rule.Field)
			}

			value = SanitizeInput(value)

			if rule.Required && value == "" {
				errors = append(errors, rule.Field+" is required")
				continue
			}

			if value == "" {
				continue
			}

			for _, r := range rule.Rules {
				switch r {
				case "email":
					if !emailRegex.MatchString(value) {
						errors = append(errors, rule.Field+" must be a valid email")
					}
				case "phone":
					if !phoneRegex.MatchString(value) {
						errors = append(errors, rule.Field+" must be a valid phone number")
					}
				case "username":
					if !usernameRegex.MatchString(value) {
						errors = append(errors, rule.Field+" must be 3-20 characters")
					}
				}
			}

			if rule.MinLen > 0 && len(value) < rule.MinLen {
				errors = append(errors, rule.Field+" must be at least "+strconv.Itoa(rule.MinLen)+" characters")
			}

			if rule.MaxLen > 0 && len(value) > rule.MaxLen {
				errors = append(errors, rule.Field+" must be at most "+strconv.Itoa(rule.MaxLen)+" characters")
			}
		}

		if len(errors) > 0 {
			c.JSON(400, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
