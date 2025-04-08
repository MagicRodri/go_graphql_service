package api

import (
	"github.com/MagicRodri/go_graphql_service/internal/db"
	"github.com/gofiber/fiber/v2"
)

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
	var result string
	err = db.GetDB().QueryRowContext(c.Context(), "SELECT graphql.resolve($1)", query).Scan(&result)
	if err != nil {
		res.Message = "Query execution failed"
		res.ResponseStatus = fiber.ErrInternalServerError.Code
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	data, total, err := transformResponse(result)
	if err != nil {
		res.Message = "Failed to transform response"
		res.ResponseStatus = fiber.ErrInternalServerError.Code
		return c.Status(fiber.StatusInternalServerError).JSON(res)
	}
	res.Data = data
	res.ResponseStatus = fiber.StatusOK
	res.Message = "Success"
	res.Count = len(data)
	res.CurrentPage = req.Page
	res.PageCount = total / req.PageSize
	res.PageSize = req.PageSize
	if len(data) == 0 {
		res.Message = "No data found"
		res.ResponseStatus = fiber.StatusNotFound
		return c.Status(fiber.StatusNotFound).JSON(res)
	}
	return c.JSON(res)
}
