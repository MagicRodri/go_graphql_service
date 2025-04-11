package api

import (
	"encoding/json"
	"fmt"

	"github.com/MagicRodri/go_graphql_service/internal/db"
	"github.com/MagicRodri/go_graphql_service/internal/logging"
	"github.com/gofiber/fiber/v2"
)

func rawGraphqlHandler(c *fiber.Ctx) error {
	var req RawRequest
	var res RawResponse
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": fmt.Sprintf("Invalid request body: %v", err),
		})
	}
	logging.Get().Debugf("GraphQL Query: %s", req.Query)

	// Execute the GraphQL query
	var result string
	if err := db.GetDB().QueryRowContext(c.Context(), "SELECT graphql.resolve($1)", req.Query).Scan(&result); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fmt.Sprintf("Query execution failed: %v", err),
		})
	}
	logging.Get().Debugf("GraphQL Response: %s", result)

	if err := json.Unmarshal([]byte(result), &res); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"errors": fmt.Sprintf("Failed to parse response: %v", err),
		})
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func graphqlHandler(c *fiber.Ctx) error {
	var req RequestDTO
	var res ResponseDTO

	if err := c.BodyParser(&req); err != nil {
		res.Message = "Invalid request body"
		res.ResponseStatus = fiber.ErrBadRequest.Code
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}

	camelTable := snakeToCamel(req.Table)
	collectionName := camelTable + "Collection"

	query, err := buildQuery(collectionName, req)
	if err != nil {
		res.Message = "Failed to build query"
		res.ResponseStatus = fiber.ErrBadRequest.Code
		return c.Status(fiber.StatusBadRequest).JSON(res)
	}
	logging.Get().Debugf("GraphQL Query: %s", query)
	var result string
	err = db.GetDB().QueryRowContext(c.Context(), "SELECT graphql.resolve($1)", query).Scan(&result)
	if err != nil {
		res.Message = "Query execution failed"
		res.ResponseStatus = fiber.ErrInternalServerError.Code
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	logging.Get().Debugf("GraphQL Response: %s", result)
	data, total, err := transformResponse(result)
	if err != nil {
		res.Message = err.Error()
		res.ResponseStatus = fiber.ErrInternalServerError.Code
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	rawResponseToDTO(&req, data, &res, total)
	return c.Status(res.ResponseStatus).JSON(res)
}

func rawResponseToDTO(req *RequestDTO, rawData []map[string]interface{}, res *ResponseDTO, total int) {

	if total == 0 {
		res.Message = "No data found"
		res.ResponseStatus = fiber.StatusNotFound
	}
	res.Data = rawData
	res.Count = total
	res.CurrentPage = req.Page
	res.PageCount = total / req.PageSize
	res.PageSize = req.PageSize
	res.ResponseStatus = fiber.StatusOK
	res.Message = "Success"
}
