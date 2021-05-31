package types

type IMerkleNode interface {
	GetHash() (hash []byte, err error)
}

type StorageItemMerkleNode StorageItem

func (self *StorageItemMerkleNode) GetHash() (hash []byte, err error) {
	item := StorageItem(*self)
	return item.GetHash()
}

type DummyMerkleNode struct {
}

func (self *DummyMerkleNode) GetHash() (hash []byte, err error) {
	return make([]byte, DefaultHashProvider.DigestSize()), nil
}
