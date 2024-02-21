package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetXKey(t *testing.T) {
	retRaw := getXKey()
	ret, ok := retRaw.(map[string]interface{})
	assert.True(t, ok)

	assert.NotEmpty(t, ret["public"])
	assert.NotEmpty(t, ret["seed"])
}
