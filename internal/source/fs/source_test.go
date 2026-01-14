package fs

import (
	"context"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"
)

func TestSource_Load_OK(t *testing.T) {
	fsys := fstest.MapFS{
		"001_init.up.sql":    {Data: []byte("CREATE TABLE test();")},
		"001_init.down.sql":  {Data: []byte("DROP TABLE test;")},
		"002_users.up.sql":   {Data: []byte("CREATE TABLE users();")},
		"002_users.down.sql": {Data: []byte("DROP TABLE users;")},
	}

	src := NewSource(fsys)

	migrations, err := src.Load(context.Background())
	require.NoError(t, err)
	require.Len(t, migrations, 2)

	require.Equal(t, "001_init", migrations[0].ID)
	require.Equal(t, "002_users", migrations[1].ID)

	require.Equal(t, "CREATE TABLE test();", string(migrations[0].Up))
	require.Equal(t, "DROP TABLE test;", string(migrations[0].Down))
}

func TestSource_Load_WithoutDown(t *testing.T) {
	fsys := fstest.MapFS{
		"001_init.up.sql": {Data: []byte("CREATE TABLE test();")},
	}

	src := NewSource(fsys)

	migrations, err := src.Load(context.Background())
	require.NoError(t, err)
	require.Len(t, migrations, 1)

	require.Equal(t, "001_init", migrations[0].ID)
	require.NotEmpty(t, migrations[0].Up)
	require.Empty(t, migrations[0].Down)
}

func TestSource_Load_MissingUp(t *testing.T) {
	fsys := fstest.MapFS{
		"001_init.down.sql": {Data: []byte("DROP TABLE test;")},
	}

	src := NewSource(fsys)

	_, err := src.Load(context.Background())
	require.Error(t, err)
	require.Contains(t, err.Error(), "has no up.sql")
}

func TestSource_Load_IgnoresUnknownFiles(t *testing.T) {
	fsys := fstest.MapFS{
		"README.md":         {Data: []byte("hello")},
		"001_init.up.sql":   {Data: []byte("up")},
		"001_init.down.sql": {Data: []byte("down")},
		"random.txt":        {Data: []byte("???")},
	}

	src := NewSource(fsys)

	migrations, err := src.Load(context.Background())
	require.NoError(t, err)
	require.Len(t, migrations, 1)
	require.Equal(t, "001_init", migrations[0].ID)
}

func TestSource_Load_SortedByID(t *testing.T) {
	fsys := fstest.MapFS{
		"010_last.up.sql":   {Data: []byte("up")},
		"002_middle.up.sql": {Data: []byte("up")},
		"001_first.up.sql":  {Data: []byte("up")},
	}

	src := NewSource(fsys)

	migrations, err := src.Load(context.Background())
	require.NoError(t, err)

	require.Equal(t, []string{
		"001_first",
		"002_middle",
		"010_last",
	}, []string{
		migrations[0].ID,
		migrations[1].ID,
		migrations[2].ID,
	})
}
