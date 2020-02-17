package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBump(t *testing.T) {

	data := "0.0.1"
	old, new, _, newcontent, err := BumpInContent([]byte(data), "p")

	assert.Nil(t, err)
	assert.Equal(t, "0.0.1", old)
	assert.Equal(t, "0.0.2", new)
	assert.Equal(t, "0.0.2", string(newcontent))

}
