package web

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"

	"qad/pkg/data"
)

// home is the handler for the root URL /
func home(c echo.Context) error {
	tables, err := data.GetDB().ListTables()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.Render(http.StatusOK, "base.html", map[string]interface{}{
		"tables": tables,
	})
}

// query is the handler for the /query endpoint.
func query(c echo.Context) error {
	query := c.QueryParam("q")
	if query == "" {
		return c.String(http.StatusBadRequest, "A query is required in the `q` parameter")
	}

	db := data.GetDB()
	JSONOutput, err := db.GetQueryJSON(query)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSONBlob(http.StatusOK, JSONOutput)
}

// ingest is the handler for the /data endpoint.
func ingest(c echo.Context) error {
	tableName := c.QueryParam("table")
	if tableName == "" {
		return c.String(http.StatusBadRequest, "Table name is required")
	}

	JSONPayload := make(map[string]interface{})
	bodyBytes, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid body: %s", err))
	}
	if len(bodyBytes) == 0 {
		return c.String(http.StatusBadRequest, "Empty body")
	}

	// Ensure we have valid JSON
	err = json.Unmarshal(bodyBytes, &JSONPayload)
	if err != nil {
		return c.String(http.StatusBadRequest, fmt.Sprintf("Invalid JSON: %s", err))
	}

	// Create a mapping of our columns to types
	colMap, err := data.CreateColumnMap(JSONPayload)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	db := data.GetDB()

	// Ensure the table can support the shape of the data
	err = db.EnsureTable(tableName, colMap)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// Insert the data
	err = db.Insert(tableName, JSONPayload)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, fmt.Sprintf(`{"table": "%s"}`, tableName))
}

// This was a neat way to jam something out and see it work
// but there is not a similar shortcut for altering the table
// file, err := os.CreateTemp("", "input.json")
// if err != nil {
// 	log.Fatal(err)
// }
// defer os.Remove(file.Name())
//
// file.Write(bodyBytes)
//
// db := data.GetDB()
//
// var stmtTable string
// if db.TableExists(tableName) {
// 	stmtTable = fmt.Sprintf("INSERT INTO %s", tableName)
// } else {
// 	stmtTable = fmt.Sprintf("CREATE TABLE %s AS", tableName)
// }
//
// JSONSelect := fmt.Sprintf(`%s SELECT * FROM read_json_auto('%s')`, stmtTable, file.Name())
// res, err := db.Connection.Exec(JSONSelect)
// if err != nil {
// 	return c.String(http.StatusInternalServerError, err.Error())
// }
// END HACK
