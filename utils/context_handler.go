package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gmtcore/infrastructure"
	"gmtcore/markers"
	"strconv"
	"strings"
)

func GetPathVariableIdFromContext(context *fiber.Ctx) uint {
	idStr := context.Params("id")
	intID, err := strconv.Atoi(idStr)
	if err != nil {
		context.Status(fiber.StatusBadRequest).JSON(BuildErrorResponse("Failed to process request", err.Error(), EmptyObj{}))
		return 0
	}
	return uint(intID)
}

func GetBodyFromContext(context *fiber.Ctx, dto markers.IDto) {
	if err := context.BodyParser(dto); err != nil {
		context.Status(fiber.StatusBadRequest).JSON(BuildErrorResponse("Failed to process request", err.Error(), EmptyObj{}))
		return
	}
	return
}

func GetUserCredentialFromContext(c *fiber.Ctx) (infrastructure.UserCredential, error) {
	userCredentialInterface := c.Locals("UserCredential")

	if userCredentialInterface == nil {
		return infrastructure.UserCredential{}, fmt.Errorf("no user credential found in context")
	}

	userCredential, ok := userCredentialInterface.(infrastructure.UserCredential)
	if !ok {
		return infrastructure.UserCredential{}, fmt.Errorf("invalid user credential found in context")
	}

	return userCredential, nil
}
func GetUserIdFromContext(c *fiber.Ctx, jwtService infrastructure.JwtService) uint {
	// user_id claim'ini context'ten alma
	userCredential, err := GetUserCredentialFromContext(c)
	if err != nil {
		c.Status(fiber.StatusForbidden).JSON(BuildErrorResponse("Invalid user credential", err.Error(), EmptyObj{}))
		return 0
	}
	userID, err := jwtService.GetValueByKeyFromJWTClaims(userCredential.Claims, "user_id")
	if err != nil {
		c.Status(fiber.StatusForbidden).JSON(BuildErrorResponse("Invalid user claims", err.Error(), EmptyObj{}))
		return 0
	}
	userIDFloat, ok := userID.(float64)
	if !ok {
		c.Status(fiber.StatusForbidden).JSON(BuildErrorResponse("Invalid or missing user_id", "user_id is not of type float64", EmptyObj{}))
		return 0
	}
	// userID'yi uint'e çevirme ve 0 kontrolü
	userUint := uint(userIDFloat)
	if userUint == 0 {
		c.Status(fiber.StatusForbidden).JSON(BuildErrorResponse("Invalid or missing user_id", "user_id is either not of type uint or is zero", EmptyObj{}))
		return 0
	}
	return userUint
}
func GetTokenFromHeader(c *fiber.Ctx) (string, error) {
	bearerToken := c.Get("Authorization")
	if strings.HasPrefix(bearerToken, "Bearer ") {
		return strings.TrimPrefix(bearerToken, "Bearer "), nil
	}
	return "", fmt.Errorf("Authorization header is missing")
}
