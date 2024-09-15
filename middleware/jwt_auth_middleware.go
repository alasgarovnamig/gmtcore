package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/alasgarovnamig/gmtcore/client"
	"github.com/gofiber/fiber/v2"

	dtos_request "github.com/alasgarovnamig/gmtcore/dto/request"
	"github.com/alasgarovnamig/gmtcore/dto/response"
	"github.com/alasgarovnamig/gmtcore/infrastructure"
	"github.com/alasgarovnamig/gmtcore/utils"
)

func AuthenticationMiddleware(jwtService infrastructure.JwtService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, err := utils.GetTokenFromHeader(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{}))
		}

		validateToken, err := jwtService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{}))
		}
		err = jwtService.TokenClaimsSetToContext(c, validateToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{}))
		}
		return c.Next()
	}
}

// AuthorizationApiCheckerMiddleware Check API Access
func AuthorizationApiCheckerMiddleware(client client.Client, jwtService infrastructure.JwtService, expectedApiPermissionID uint) fiber.Handler {

	return func(c *fiber.Ctx) error {
		opaRequestBody := &dtos_request.OPAAPIAuthorizationRequestDto{
			Input: dtos_request.OPAAPIAuthorizationRequestInput{
				UserID:        utils.GetUserIdFromContext(c, jwtService),
				ApiPermission: expectedApiPermissionID,
			},
		}
		resp, httpStatusCode, err := client.Post("https://", opaRequestBody, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			//TODO: uygun cavab don
			return err
		}
		//// JSON formatına çevirme
		//jsonData, err := json.Marshal(opaRequestBody)
		//if err != nil {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware JSON Serialize Error").Error(), utils.EmptyObj{}))
		//}
		//req, err := http.NewRequest("POST", config.Conf.Application.OPA.APIAuthorizationUrl, bytes.NewBuffer(jsonData))
		//if err != nil {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Error").Error(), utils.EmptyObj{}))
		//}
		//
		//req.Header.Set("Content-Type", "application/json")
		//
		//// HTTP Client ile isteği gönderme
		//client := &http.Client{}
		//resp, err := client.Do(req)
		//if err != nil || resp.StatusCode != 200 {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Access Manager Or Network Error").Error(), utils.EmptyObj{}))
		//}
		//defer resp.Body.Close()
		//
		//body, err := io.ReadAll(resp.Body)
		//if err != nil {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Access Manager Read Response Body Error").Error(), utils.EmptyObj{}))
		//}

		var responseBody response.OPAAPIAuthorizationResponseDto
		err = json.Unmarshal(resp, &responseBody)
		if err != nil || httpStatusCode != 200 {
			return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Access Manager Read Response Body Casting Error").Error(), utils.EmptyObj{}))
		}

		if !responseBody.Result.Allow {
			return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("You do not have permission to access this resource", "You do not have permission to access this resource", utils.EmptyObj{}))
		}
		userCredential, err := utils.GetUserCredentialFromContext(c)
		userCredential.UserParentID = responseBody.Result.ParentUserID
		userCredential.UserRoleID = responseBody.Result.RoleID
		//c.Locals("role_id", responseBody.Result.RoleID)
		//c.Locals("parent_user_id", responseBody.Result.ParentUserID)
		c.Locals("UserCredential", userCredential)
		return c.Next()
	}
}

// AuthorizationSearchFieldCheckerMiddleware Check Filed Access For Search Operation
func AuthorizationSearchFieldCheckerMiddleware(client client.Client, jwtService infrastructure.JwtService, sourceTable string, readableTables []string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		searchRequestDto := &dtos_request.SearchRequestDto{}
		if err := c.BodyParser(searchRequestDto); err != nil {
			return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Invalid or missing user_id", "user_id is either not of type uint or is zero", utils.EmptyObj{}))
		}
		opaRequestBody := &dtos_request.OPASearchFieldCheckerRequestDto{
			Input: dtos_request.OPASearchFieldCheckerInput{
				UserID:          utils.GetUserIdFromContext(c, jwtService),
				SearchableTable: sourceTable,
				ReadTables:      readableTables,
				UserPermissions: make([]dtos_request.UserPermission, len(searchRequestDto.Criteria)),
			},
		}
		for index, criterion := range searchRequestDto.Criteria {
			perm := dtos_request.UserPermission{
				Key:             criterion.Key,
				FieldPermission: 4,
				SearchCriteria:  uint(criterion.Operation),
			}
			opaRequestBody.Input.UserPermissions[index] = perm
		}

		resp, httpStatusCode, err := client.Post("https://", opaRequestBody, map[string]string{
			"Content-Type": "application/json",
		})
		if err != nil {
			//TODO: uygun cavab don
			return err
		}
		//jsonData, err := json.Marshal(opaRequestBody)
		//if err != nil {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware JSON Serialize Error").Error(), utils.EmptyObj{}))
		//}
		//
		//req, err := http.NewRequest("POST", "http://localhost:8181/v1/data/authz/search_field/result", bytes.NewBuffer(jsonData))
		//if err != nil {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Error").Error(), utils.EmptyObj{}))
		//}
		//
		//req.Header.Set("Content-Type", "application/json")
		//
		//// HTTP Client ile isteği gönderme
		//client := &http.Client{}
		//resp, err := client.Do(req)
		//if err != nil || resp.StatusCode != 200 {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Access Manager Or Network Error").Error(), utils.EmptyObj{}))
		//}
		//defer resp.Body.Close()
		//
		//body, err := io.ReadAll(resp.Body)
		//if err != nil || resp.StatusCode != 200 {
		//	return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Access Manager Read Response Body Error").Error(), utils.EmptyObj{}))
		//}

		var responseBody response.OPASearchFieldAuthorizationResponseDto
		err = json.Unmarshal(resp, &responseBody)
		if err != nil || httpStatusCode != 200 {
			return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("Failed to process request", fmt.Errorf("AuthorizationRBACMiddleware Check Api Access Request Access Manager Read Response Body Casting Error").Error(), utils.EmptyObj{}))
		}

		if !responseBody.Result.Allow {
			return c.Status(fiber.StatusForbidden).JSON(utils.BuildErrorResponse("You do not have permission to access this resource", "You do not have permission to access this resource", utils.EmptyObj{}))
		}
		c.Locals("readable_fields", responseBody.Result.ReadableFields)
		c.Locals("readable_tables", readableTables)
		return c.Next()
	}
}

// AuthorizationGetByIDFieldCheckerMiddleware Check Filed Access For Get By ID Operation
func AuthorizationGetByIDFieldCheckerMiddleware(c *fiber.Ctx) fiber.Handler {

	return func(c *fiber.Ctx) error {

		return c.Next()
	}
}

// AuthorizationCreateFieldCheckerMiddleware Check Filed Access For Create Operation
func AuthorizationCreateFieldCheckerMiddleware(c *fiber.Ctx) fiber.Handler {

	return func(c *fiber.Ctx) error {

		return c.Next()
	}
}

// AuthorizationUpdateFieldCheckerMiddleware Check Filed Access For Update Operation
func AuthorizationUpdateFieldCheckerMiddleware(c *fiber.Ctx) fiber.Handler {

	return func(c *fiber.Ctx) error {

		return c.Next()
	}
}
