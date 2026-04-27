package fs

import (
	"io/fs"

	internalfs "github.com/Palladium-blockchain/go-migrations/internal/source/fs"
)

type Source = internalfs.Source

func NewSource(fsys fs.FS) *Source {
	return internalfs.NewSource(fsys)
}
