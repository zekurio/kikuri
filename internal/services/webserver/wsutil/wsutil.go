package wsutil

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetQueryInt tries to get a value from request query
// and returns the value as int.
//
// If the value is not found, the default value is returned.
//
// If the value is smaller than min or bigger than max,
// 0 and an error is returned.
func GetQueryInt(ctx *fiber.Ctx, key string, defaultValue, min, max int) (int, error) {
	value := ctx.Query(key)
	if value == "" {
		return defaultValue, nil
	}

	res, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	if res < min || res > max {
		return 0, fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("value of '%s' must be in bounds [%d, %d]", key, min, max))
	}

	return res, nil
}

// GetQueryBool tries to get a value from request query
// and returns the value as bool.
//
// If the value is not found, the default value is returned.
func GetQueryBool(ctx *fiber.Ctx, key string, defaultValue bool) (bool, error) {
	value := ctx.Query(key)
	if value == "" {
		return defaultValue, nil
	}

	switch strings.ToLower(value) {
	case "true", "1", "yes", "y":
		return true, nil
	case "false", "0", "no", "n":
		return false, nil
	default:
		return false, fiber.NewError(fiber.StatusBadRequest,
			fmt.Sprintf("value of '%s' must be a boolean", key))
	}
}
