package storage

type StorageWrapper interface {
	Save([]byte) error
}
