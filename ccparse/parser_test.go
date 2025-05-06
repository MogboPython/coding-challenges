package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckValidity(t *testing.T) {
	t.Run("Check if json starts with { and ends with }", func(t *testing.T) {

		dummyContent := "{}"
		r, err := check_validity([]byte(dummyContent))

		assert.NotNil(t, r)
		assert.Nil(t, err)
	})

	t.Run("Check if json is valid and parses {} and []", func(t *testing.T) {

		dummyContent := "{\"key\": \"value\",\"key-n\": 101,\"key-o\": {},\"key-l\": []}"
		r, err := check_validity([]byte(dummyContent))

		assert.NotNil(t, r)
		assert.Nil(t, err)
	})

	t.Run("Check if json starts with { and ends with", func(t *testing.T) {

		dummyContent := "{"
		r, err := check_validity([]byte(dummyContent))

		assert.Nil(t, r)
		assert.Equal(t, err, errors.New("JSON string does not end with }"))
	})

	t.Run("Check if json starts contains only }", func(t *testing.T) {

		dummyContent := "}"
		r, err := check_validity([]byte(dummyContent))

		assert.Nil(t, r)
		assert.Equal(t, err, errors.New("JSON string does not start with {"))
	})

	t.Run("Check if json is invalid", func(t *testing.T) {

		dummyContent := ""
		r, err := check_validity([]byte(dummyContent))
		assert.Nil(t, r)
		assert.Equal(t, err, errors.New("empty JSON string"))
	})

	t.Run("Check invalid json that ends with : ", func(t *testing.T) {

		dummyContent := "{\"key\":}"
		r, err := check_validity([]byte(dummyContent))

		assert.Equal(t, err, errors.New("invalid JSON: colon must be followed by a value"))
		assert.Nil(t, r)
	})
	t.Run("Check invalid json that with non string key ", func(t *testing.T) {

		dummyContent := "{\"key\": \"value\", key2: \"value\"}"
		r, err := check_validity([]byte(dummyContent))

		assert.Equal(t, err, errors.New("invalid JSON: key must be a quoted string, found \"key2\""))
		assert.Nil(t, r)
	})

	t.Run("Check invalid json that with invalid value", func(t *testing.T) {

		dummyContent := "{\"key\": \"value\", \"key2\": False}"
		r, err := check_validity([]byte(dummyContent))

		assert.Equal(t, err, errors.New("invalid JSON: value unexpected"))
		assert.Nil(t, r)
	})

	t.Run("Check invalid json that ends with ,", func(t *testing.T) {

		dummyContent := "{\"key\": \"value\",}"
		r, err := check_validity([]byte(dummyContent))

		assert.Equal(t, err, errors.New("invalid JSON: comma must be followed by a field or object end"))
		assert.Nil(t, r)
	})
}
