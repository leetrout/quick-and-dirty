package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"qad/pkg/data"
)

func TestIngest(t *testing.T) {
	sampleTable := `{
		"foo": 1,
		"bar": "baz"
	}`

	db := data.Connect(data.ConnectWithDSN("test.duckdb"))
	defer db.Connection.Close()

	e := NewEcho()

	url := e.Reverse("ingest") + "?table=foobar"
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(sampleTable))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedResp := `{"table": "foobar"}`

	if assert.NoError(t, ingest(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	}
}

func TestQuery(t *testing.T) {
	db := data.Connect(data.ConnectWithDSN("test.duckdb"))
	defer db.Connection.Close()

	e := NewEcho()

	url := e.Reverse("query") + "?q=SELECT+*+FROM+foobar"
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedResp := `[{"bar":"baz","foo":1}]`

	if assert.NoError(t, query(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	}

	t.Cleanup(func() {
		os.Remove("test.duckdb")
	})
}
