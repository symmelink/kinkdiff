package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestInit(t *testing.T) {
	t.Log("init ran!")

	require.NotEmpty(t, Quiz, "quiz categories")
	require.NotEmpty(t, Quiz[0].Title, "quiz[0] title")
}
