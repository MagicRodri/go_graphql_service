package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/MagicRodri/go_graphql_service/internal/logging"
	"github.com/gofiber/fiber"
)

func buildQuery(collectionName string, req RequestDTO) (string, error) {
	var sb strings.Builder
	sb.WriteString("query { ")
	sb.WriteString(fmt.Sprintf("%s(", collectionName))

	if req.PageSize > 0 {
		sb.WriteString(fmt.Sprintf("first: %d, ", req.PageSize))
	}
	req.Page -= 1
	if req.Page > 0 && req.PageSize > 0 {
		sb.WriteString(fmt.Sprintf("offset: %d, ", req.Page*req.PageSize))
	}
	logging.Get().Debugf("Filters: %v, len: %d", req.Filters, len(req.Filters))
	if len(req.Filters) > 0 {
		filter, err := serializeFilter(req.Filters)
		if err != nil {
			logging.Get().Debugf("Error serializing filter: %v", err)
			return "", err
		}
		sb.WriteString(fmt.Sprintf("filter: %s, ", filter))
	}

	queryPart := strings.TrimSuffix(sb.String(), ", ") + ") {"
	sb.Reset()
	sb.WriteString(queryPart)

	sb.WriteString("edges { node { ")
	if len(req.Fields) == 0 {
		sb.WriteString("id ")
	} else {
		for _, field := range req.Fields {
			sb.WriteString(fmt.Sprintf("%s ", field))
		}
	}

	// Relationships
	for relName, rel := range req.Extra {
		if strings.HasSuffix(relName, "_set") {
			// One-to-many relationship
			collection := strings.TrimSuffix(relName, "_set") + "Collection"
			if err := addCollection(&sb, snakeToCamel(collection), rel); err != nil {
				return "", err
			}
		} else {
			// Many-to-one relationship
			if err := addNode(&sb, relName, rel); err != nil {
				return "", err
			}
		}
	}

	sb.WriteString("} } pageInfo { hasNextPage hasPreviousPage }")
	sb.WriteString(" totalCount")
	sb.WriteString(" } }")
	return sb.String(), nil
}

func addCollection(sb *strings.Builder, name string, rel Relation) error {
	addParams(sb, rel.PageSize, rel.Page, rel.Filters)
	sb.WriteString(fmt.Sprintf("%s { edges { node { ", name))

	// Collection fields
	if len(rel.Fields) == 0 {
		sb.WriteString("id ")
	} else {
		for _, field := range rel.Fields {
			sb.WriteString(fmt.Sprintf("%s ", snakeToCamel(field)))
		}
	}

	// Nested relationships
	for nestedName, nestedRel := range rel.Extra {
		if strings.HasSuffix(nestedName, "_set") {
			nestedCollection := strings.TrimSuffix(nestedName, "_set") + "Collection"
			if err := addCollection(sb, snakeToCamel(nestedCollection), nestedRel); err != nil {
				return err
			}
		} else {
			if err := addNode(sb, nestedName, nestedRel); err != nil {
				return err
			}
		}
	}

	if rel.TotalCount {
		sb.WriteString("totalCount ")
	}

	sb.WriteString("} } } ")
	return nil
}

func addNode(sb *strings.Builder, name string, rel Relation) error {
	sb.WriteString(fmt.Sprintf("%s { ", name))

	// Node fields
	if len(rel.Fields) == 0 {
		sb.WriteString("id ")
	} else {
		for _, field := range rel.Fields {
			sb.WriteString(fmt.Sprintf("%s ", field))
		}
	}

	// Nested relationships
	for nestedName, nestedRel := range rel.Extra {
		if strings.HasSuffix(nestedName, "_set") {
			nestedCollection := strings.TrimSuffix(nestedName, "_set") + "Collection"
			if err := addCollection(sb, snakeToCamel(nestedCollection), nestedRel); err != nil {
				return err
			}
		} else {
			if err := addNode(sb, nestedName, nestedRel); err != nil {
				return err
			}
		}
	}
	sb.WriteString("} ")
	return nil
}

func addParams(sb *strings.Builder, pageSize, page int, filters map[string]interface{}) {
	if pageSize > 0 {
		sb.WriteString(fmt.Sprintf("    first: %d,\n", pageSize))
	}
	if page > 0 {
		sb.WriteString(fmt.Sprintf("offset: %d, ", (page)*pageSize))
	}
	if len(filters) > 0 {
		filter, err := serializeFilter(filters)
		if err == nil {
			sb.WriteString(fmt.Sprintf("filter: %s, ", filter))
		}
	}
}

func transformResponse(raw string) ([]map[string]interface{}, int, error) {
	var res map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return nil, 0, err
	}

	var data []map[string]interface{}
	var total int

	if errList, ok := res["errors"].([]interface{}); ok {
		for _, err := range errList {
			if errMap, ok := err.(map[string]interface{}); ok {
				if message, ok := errMap["message"].(string); ok {
					return nil, 0, fmt.Errorf("error: %s", message)
				}
			}
		}
	}

	if dataMap, ok := res["data"].(map[string]interface{}); ok {
		for _, collection := range dataMap {
			if coll, ok := collection.(map[string]interface{}); ok {
				if edges, ok := coll["edges"].([]interface{}); ok {
					for _, e := range edges {
						if edge, ok := e.(map[string]interface{}); ok {
							if node, ok := edge["node"].(map[string]interface{}); ok {
								data = append(data, node)
							}
						}
					}
				}
				if totalCount, ok := coll["totalCount"].(float64); ok {
					total = int(totalCount)
				}
			}
		}
	}

	return data, total, nil
}

func snakeToCamel(s string) string {
	words := strings.Split(s, "_")
	for i := 1; i < len(words); i++ {
		words[i] = string(unicode.ToUpper(rune(words[i][0]))) + words[i][1:]
	}
	return strings.Join(words, "")
}

func serializeFilter(filters map[string]interface{}) (string, error) {
	var parts []string
	for field, condition := range filters {
		condMap, ok := condition.(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid condition for field %s", field)
		}

		var condParts []string
		for op, value := range condMap {
			valueJSON, err := json.Marshal(value)
			if err != nil {
				return "", fmt.Errorf("invalid value for %s.%s", field, op)
			}
			condParts = append(condParts, fmt.Sprintf("%s: %s", op, string(valueJSON)))
		}
		parts = append(parts, fmt.Sprintf("%s: { %s }", field, strings.Join(condParts, ", ")))
	}
	return fmt.Sprintf("{ %s }", strings.Join(parts, ", ")), nil
}

func InitGraphQLTables(db *sql.DB) error {
	// Get all table names from public schema
	rows, err := db.Query(`
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'public' 
		AND table_type = 'BASE TABLE'
	`)
	if err != nil {
		return fmt.Errorf("failed to query tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %v", err)
		}
		tables = append(tables, tableName)
	}

	for _, table := range tables {
		if err := enableTotalCount(db, table); err != nil {
			return fmt.Errorf("failed to enable totalCount for %s: %v", table, err)
		}
		logging.Get().Debug("Enabled totalCount for table: %s", table)
	}

	return nil
}

func enableTotalCount(db *sql.DB, tableName string) error {
	// Check if totalCount is already enabled
	var existingComment sql.NullString
	err := db.QueryRow(`
		SELECT obj_description($1::regclass)
	`, tableName).Scan(&existingComment)
	if err != nil {
		return fmt.Errorf("failed to check existing comment: %v", err)
	}

	// Skip if already configured
	if strings.Contains(existingComment.String, `"totalCount": {"enabled": true}`) {
		return nil
	}

	_, err = db.Exec(fmt.Sprintf(`
		COMMENT ON TABLE "%s" IS e'@graphql({"totalCount": {"enabled": true}})'
	`, tableName))
	if err != nil {
		return fmt.Errorf("failed to set comment: %v", err)
	}

	return nil
}

func rawResponseToDTO(req *RequestDTO, rawData []map[string]interface{}, res *ResponseDTO, total int) {

	res.ResponseStatus = fiber.StatusOK
	if total == 0 {
		res.Message = "No data found"
		res.ResponseStatus = fiber.StatusNotFound
	}
	res.Data = rawData
	res.Count = total
	res.CurrentPage = req.Page
	res.PageCount = int(math.Ceil(float64(total) / float64(req.PageSize)))
	res.PageSize = len(rawData)
}
