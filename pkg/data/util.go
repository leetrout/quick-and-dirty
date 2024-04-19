package data

import (
	"fmt"
	"reflect"
	"regexp"
)

const TABLE_NAME_PATTERN = `^[a-zA-Z]\w{2,47}$`

var rxTableName = regexp.MustCompile(TABLE_NAME_PATTERN)

// ValidateTableName ensures the table name:
// - starts with an ascii letter
// - is between 3 and 48 chars in length
// - only contains word characters
// TODO: More human friendly / detailed error message
func ValidateTableName(tableName string) error {
	if !rxTableName.MatchString(tableName) {
		return fmt.Errorf("invalid table name `%s`", tableName)
	}
	return nil
}

// ColumnMap represents a DB table.
// cm["name"] = "VARCHAR"
type ColumnMap map[string]string

// AsTuple returns the column map as a tuple for use in table create statements.
func (cm ColumnMap) AsTuple() string {
	cols := "("
	idx := 0
	for col, kind := range cm {
		if idx > 0 {
			cols += ", "
		}
		cols += fmt.Sprintf(`"%s" %s`, col, kind)
		idx++
	}
	cols += ")"
	return cols
}

// CreateColumnMap takes a JSON payload and converts the keys into column names
// and duckdb column types.
// TODO: Needs col name validation and to return err
func CreateColumnMap(JSONPayload map[string]interface{}) (ColumnMap, error) {
	cm := ColumnMap{}
	var colType string
	for key, ival := range JSONPayload {
		kind := reflect.TypeOf(ival)

		switch ival.(type) {
		case float32:
			colType = "FLOAT"
		case float64:
			colType = "DOUBLE"
		case int:
			colType = "BIGINT"
		case int32:
			colType = "BIGINT"
		case int64:
			colType = "BIGINT"
		case string:
			colType = "VARCHAR"
		case bool:
			colType = "BOOLEAN"
		default:
			colType = ""
		}

		if colType == "" {
			return nil, fmt.Errorf("field `%s` has unsupported type `%s`", key, kind)
		}
		cm[key] = colType
	}
	return cm, nil
}
