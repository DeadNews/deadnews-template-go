package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnswer(t *testing.T) {
	assert := assert.New(t)

	assert.Equal(42, GetAnswer(), "they should be equal")
}
