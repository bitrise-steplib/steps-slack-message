package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ensureNewlineEscapeChar(t *testing.T) {
	require.Equal(t, "", ensureNewlineEscapeChar(""))
	require.Equal(t, "a", ensureNewlineEscapeChar("a"))
	require.Equal(t, "\n", ensureNewlineEscapeChar("\n"))
	require.Equal(t, "\n", ensureNewlineEscapeChar("\\"+"n"))
	// should convert \ + n to \n too; where \n is a single char, the ASCII 10 "newline feed" char
	require.Equal(t, uint8(10), ensureNewlineEscapeChar("\n")[0])
	require.Equal(t, uint8(10), ensureNewlineEscapeChar("\\" + "n")[0])
}
