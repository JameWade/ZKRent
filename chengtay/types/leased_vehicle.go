package types

type Vehicle struct {
	VehicleOwnerId   ID                         `json:"carOwnerId"`
	VehicleId ID								`json:"carId"`
	UsedTime    string						`json:"usedTime"`
}

func (self *Vehicle) GetId() ID {
	return self.VehicleId
}

func (self *Vehicle) GetCarOwnerId() ID {
	return self.VehicleOwnerId
}

//compute clearing every week or month?
func (self *Vehicle) GetUseTime() string {
	return self.UsedTime
}

//func(self *Vehicle) Marshal() []byte {
//	var buffer bytes.Buffer
//	buffer.Write([]byte(self.CarOwnerId))
//	buffer.Write([]byte(self.CarId))
//	buffer.Write([]byte(self.UsedTime))
//	return buffer.Bytes()
//}
