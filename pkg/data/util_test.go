package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTableName(t *testing.T) {
	assert.NoError(t, ValidateTableName("foo"))
	assert.NoError(t, ValidateTableName("Foob"))
	assert.NoError(t, ValidateTableName("foobar"))
	assert.NoError(t, ValidateTableName("foo_bar"))
	assert.NoError(t, ValidateTableName("fourfourfourfourfourfourfourfourfourfourfourfour"))

	assert.Error(t, ValidateTableName("1foobar"))
	assert.Error(t, ValidateTableName("foo-bar"))
	assert.Error(t, ValidateTableName("foo.bar"))
	assert.Error(t, ValidateTableName("fo"))
	assert.Error(t, ValidateTableName("fourfourfourfourfourfourfourfourfourfourfourfoura"))
}
