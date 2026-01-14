package fs

import (
	"context"
	"fmt"
	"io/fs"
	"palladium-intelligence/go-migrations/pkg/migrate"
	"sort"
	"strings"
)

type Source struct {
	fsys fs.FS
}

func NewSource(fsys fs.FS) *Source {
	return &Source{fsys: fsys}
}

func (s *Source) Load(_ context.Context) ([]migrate.Migration, error) {
	entries, err := fs.ReadDir(s.fsys, ".")
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	migrations := make(map[string]*migrate.Migration)

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		info, ok := parseFilename(e.Name())
		if !ok {
			continue
		}

		data, err := fs.ReadFile(s.fsys, e.Name())
		if err != nil {
			return nil, err
		}

		m, exists := migrations[info.ID]
		if !exists {
			m = &migrate.Migration{ID: info.ID}
			migrations[info.ID] = m
		}

		switch info.Direction {
		case "up":
			m.Up = data
		case "down":
			m.Down = data
		}
	}

	result := make([]migrate.Migration, 0, len(migrations))
	for _, m := range migrations {
		if len(m.Up) == 0 {
			return nil, fmt.Errorf("migration %s has no up.sql", m.ID)
		}
		result = append(result, *m)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result, nil
}

type fileInfo struct {
	ID        string
	Direction string
}

func parseFilename(name string) (fileInfo, bool) {
	if !strings.HasSuffix(name, ".sql") {
		return fileInfo{}, false
	}

	base := strings.TrimSuffix(name, ".sql")
	parts := strings.Split(base, ".")
	if len(parts) != 2 {
		return fileInfo{}, false
	}

	id := parts[0]
	dir := parts[1]

	if dir != "up" && dir != "down" {
		return fileInfo{}, false
	}

	return fileInfo{
		ID:        id,
		Direction: dir,
	}, true
}
