package fs

import internalfs "github.com/Palladium-blockchain/go-migrations/internal/creator/fs"

type Creator = internalfs.Creator

func NewCreator(dir string) *Creator {
	return internalfs.NewCreator(dir)
}
