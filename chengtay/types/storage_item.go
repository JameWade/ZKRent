package types

import (
	"github.com/ChengtayChain/ChengtayChain/chengtay/util"
)

type StorageItem struct {
	StorageItemID ID
	CarID         ID
	Timestamp     uint64
	ContentType   string
	Content       []byte
}

func (self *StorageItem) GetHash() (digest []byte, err error) {
	source := make([]byte, 0)

	// 1. StorageItemID
	source = append(source, []byte(self.StorageItemID)...)

	// 2. CarID
	source = append(source, []byte(self.CarID)...)

	// 3. Timestamp (uint64) in BigEndian
	{
		ret, err := util.UInt64ToBytes(self.Timestamp)
		if err != nil {
			return nil, err
		}
		source = append(source, ret...)
	}

	// 4. ContentType
	source = append(source, []byte(self.ContentType)...)

	// 5. Content
	source = append(source, []byte(self.Content)...)

	return DefaultHashProvider.Digest(source), nil
}
