package fs

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreator_Create_OK(t *testing.T) {
	dir := t.TempDir()

	c := &Creator{dir: dir}

	files, err := c.Create(context.Background(), "add_users")
	require.NoError(t, err)

	// пути должны быть внутри dir
	require.True(t, strings.HasPrefix(files.Up, dir))
	require.True(t, strings.HasPrefix(files.Down, dir))

	// файлы должны существовать
	upData, err := os.ReadFile(files.Up)
	require.NoError(t, err)
	require.Equal(t, "-- up\n", string(upData))

	downData, err := os.ReadFile(files.Down)
	require.NoError(t, err)
	require.Equal(t, "-- down\n", string(downData))
}
