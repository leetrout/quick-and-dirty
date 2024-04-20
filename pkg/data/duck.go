package data

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/marcboeker/go-duckdb"
)

type connectParam func(c *connectConfig)

type connectConfig struct {
	dsn string
}

func ConnectWithDSN(dsn string) connectParam {
	return func(c *connectConfig) {
		c.dsn = dsn
	}
}

var activeDB *DB

func Connect(connectParams ...connectParam) *DB {
	c := &connectConfig{}
	for _, cp := range connectParams {
		cp(c)
	}

	db, err := sql.Open("duckdb", c.dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Configure DB
	_, err = db.Exec("INSTALL json")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("LOAD json")
	if err != nil {
		log.Fatal(err)
	}

	// Memoize
	activeDB = &DB{
		Connection: db,
	}
	return activeDB
}

// GetDB returns the current DB connection or panics.
func GetDB() *DB {
	if activeDB == nil {
		panic("No active DB connection. Use `data.Connect` first.")
	}
	return activeDB
}

type DB struct {
	Connection *sql.DB
}

func (db *DB) Insert(tableName string, JSONPayload map[string]interface{}) error {
	params := make([]interface{}, len(JSONPayload))
	cols := "("
	vals := "("
	idx := 0
	for col, ival := range JSONPayload {
		if idx > 0 {
			vals += ", "
			cols += ", "
		}
		cols += fmt.Sprintf(`"%s"`, col)
		vals += "?"
		params[idx] = ival
		idx++
	}
	cols += ")"
	vals += ")"

	query := fmt.Sprintf("INSERT INTO %s %s VALUES %s", tableName, cols, vals)
	_, err := db.Connection.Exec(query, params...)
	return err
}

// ListTables returns the list of table names.
func (db *DB) ListTables() ([]string, error) {
	tableNames := []string{}
	rows, err := db.Connection.Query(
		"SELECT table_name FROM information_schema.tables",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			return nil, err
		}

		tableNames = append(tableNames, tableName)
	}

	return tableNames, nil
}

// TableExists checks for a table with the given name.
func (db *DB) TableExists(tableName string) bool {
	var foundTableName string
	err := db.Connection.QueryRow(
		"SELECT table_name FROM information_schema.tables WHERE table_name = ?",
		tableName,
	).Scan(&foundTableName)
	return err == nil
}

// EnsureTable creates the given table if it doesn't exist.
func (db *DB) CreateTable(tableName string, colMap ColumnMap) error {
	err := ValidateTableName(tableName)
	if err != nil {
		return err
	}

	createStatement := fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s %s",
		tableName,
		colMap.AsTuple(),
	)

	_, err = db.Connection.Exec(createStatement)
	return err
}

// EnsureTable creates the given table if it doesn't exist.
func (db *DB) EnsureTable(tableName string, colMap ColumnMap) error {
	if !db.TableExists(tableName) {
		return db.CreateTable(tableName, colMap)
	}
	return db.EnsureColumns(tableName, colMap)
}

func (db *DB) EnsureColumns(tableName string, colMap ColumnMap) error {
	// Get existing columns
	existingCols := map[string]struct{}{}

	rows, err := db.Connection.Query("select column_name from information_schema.columns where table_name = ?", tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var colName string
		err := rows.Scan(&colName)
		if err != nil {
			return err
		}
		existingCols[colName] = struct{}{}
	}

	// Diff with map
	for col, kind := range colMap {
		if _, ok := existingCols[col]; !ok {
			_, err := db.Connection.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", tableName, col, kind))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetQueryJSON returns a JSON byte slice of the results of the given
// query. Null values are not included in the output.
func (db *DB) GetQueryJSON(query string) ([]byte, error) {
	// https://stackoverflow.com/questions/42774467/how-to-convert-sql-rows-to-typed-json-in-golang/60386531#60386531
	rows, err := db.Connection.Query(query)
	if err != nil && err != sql.ErrNoRows {
		return []byte{}, err
	}
	defer rows.Close()

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return []byte{}, err
	}

	colCount := len(columnTypes)
	output := []interface{}{}

	for rows.Next() {

		// Create nullable value holders for this row
		rowVals := make([]interface{}, colCount)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR":
				rowVals[i] = new(sql.NullString)
			case "BOOLEAN":
				rowVals[i] = new(sql.NullBool)
			case "BIGINT":
				rowVals[i] = new(sql.NullInt64)
			case "DOUBLE":
				rowVals[i] = new(sql.NullFloat64)
			default:
				rowVals[i] = new(sql.NullString)
			}
		}

		err := rows.Scan(rowVals...)
		if err != nil {
			return []byte{}, err
		}

		thisRowData := map[string]interface{}{}

		// Walk each col in this row and see if the value is null (!valid)
		// and include non-null values in our final output.
		for i, v := range columnTypes {
			if z, ok := (rowVals[i]).(*sql.NullBool); ok {
				if !z.Valid {
					continue
				}
				thisRowData[v.Name()] = z.Bool
				continue
			}

			if z, ok := (rowVals[i]).(*sql.NullString); ok {
				if !z.Valid {
					continue
				}
				thisRowData[v.Name()] = z.String
				continue
			}

			if z, ok := (rowVals[i]).(*sql.NullInt64); ok {
				if !z.Valid {
					continue
				}
				thisRowData[v.Name()] = z.Int64
				continue
			}

			if z, ok := (rowVals[i]).(*sql.NullFloat64); ok {
				if !z.Valid {
					continue
				}
				thisRowData[v.Name()] = z.Float64
				continue
			}

			if z, ok := (rowVals[i]).(*sql.NullInt32); ok {
				if !z.Valid {
					continue
				}
				thisRowData[v.Name()] = z.Int32
				continue
			}

			// Default to the value present
			thisRowData[v.Name()] = rowVals[i]
		}

		output = append(output, thisRowData)
	}

	return json.Marshal(output)
}
