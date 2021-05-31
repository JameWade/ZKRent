package types

type VehicleOwnerIncomeInfo struct {
	VehicleOwnerId    ID         `json:"vehicleownerId"`
	Vehicles []*Vehicle    		 `json:"vehicle"`
	Income     string     		 `json:"income"`
	Month      int       		 `json:"month"` //按月计算   1 means Jan. and 12 means Dec.
}

func (self *VehicleOwnerIncomeInfo) GetId() ID {
	return self.VehicleOwnerId
}

func (self *VehicleOwnerIncomeInfo) GetIncome() string {
	return self.Income
}

//获取本人在平台的所有注册车辆并返回
func (self *VehicleOwnerIncomeInfo) GetOwnerCar() {
	return
}

/*
func(self *VehicleOwnerIncomeInfo) Marshal() []byte {
	var buffer bytes.Buffer
	buffer.Write([]byte(self.OwnerId))
	for i :=0 ; i < len(self.LeasedCars) ; i++{
		carBytes := self.LeasedCars[i].Marshal()
		buffer.Write(carBytes)
	}
	buffer.Write([]byte(self.Income))
	buffer.Write([]byte(self.Month.String()))
	return buffer.Bytes()
}*/
