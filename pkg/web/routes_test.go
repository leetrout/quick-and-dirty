package web

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteMap(t *testing.T) {
	name := "foo"
	url := "/foo"
	rm := RouteMap{
		name: url,
	}

	// Basic mapping exists
	assert.Equal(t, rm[name], url)
	// Nothing cached
	assert.Equal(t, len(_reversedRouteMap), 0)
	// Reverse lookup works
	assert.Equal(t, rm.Reverse(url), name)
	// Reverse map is cached
	assert.NotEqual(t, len(_reversedRouteMap), 0)
}
