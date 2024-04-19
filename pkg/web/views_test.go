package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestIngest(t *testing.T) {
	sampleTable := `{
		"foo": 1,
		"bar": "baz"
	}`

	e := echo.New()
	url := RouteLookup["ingest"] + "?table=foobar"
	req := httptest.NewRequest(http.MethodPost, url, strings.NewReader(sampleTable))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedResp := "table: foobar"

	if assert.NoError(t, ingest(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	}
}

func TestQuery(t *testing.T) {
	e := echo.New()
	url := RouteLookup["query"] + "?q=SELECT+*+FROM+foobar"
	req := httptest.NewRequest(http.MethodGet, url, strings.NewReader(""))
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	expectedResp := "query: SELECT * FROM foobar"

	if assert.NoError(t, query(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expectedResp, rec.Body.String())
	}
}
