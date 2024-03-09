package bookmarks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionsBuilder(t *testing.T) {
	builder := NewConditionBuilder()
	builder.Add("foo = bar")
	builder.Add("list IN (?, ?, ?)")
	builder.Add("id = 1452")

	expected := "WHERE foo = bar AND list IN (?, ?, ?) AND id = 1452"
	actual := builder.String()

	assert.Equal(t, expected, actual)
}
