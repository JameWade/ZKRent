package internal

import (
	"encoding/json"
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"time"
)

func GenTestVehicleData() (TestFile string){
	vehicle := []types.Vehicle{
		{"1", "0","1"},
		{"1", "1","2"},
		{"1","2","3"},
	}
	leasedCarsBytes, err := json.Marshal(vehicle)
	if err != nil{
		panic(err)
	}
	homeDir, err := homedir.Dir()
	TestFile = homeDir + string(os.PathSeparator) + "carinfo.json"
	err = ioutil.WriteFile(TestFile,leasedCarsBytes,0666)
	return TestFile
}

func GetVehicleOwner(id types.ID) (vehicleOwnerIncomeInfo *types.VehicleOwnerIncomeInfo) {
	var filename = GenTestVehicleData()
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var vehicles []types.Vehicle
	vehicles, err = UnmarshalVehicles(bytes)
	vehicleOwnerIncomeInfo = ProcessVehicles(vehicles)
	vehicleOwnerIncomeInfo.VehicleOwnerId = id
	now := time.Now()
	/*
		1. type-cast is used to convert month to int, as "time" package defines.
		2. 1 means Jan. and 12 means Dec.
	 */
	vehicleOwnerIncomeInfo.Month = int(now.Month())
	return vehicleOwnerIncomeInfo
}

func printAny(obj interface{}) {
	fmt.Println("Print")
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
	fmt.Println()
}

func UnmarshalVehicles(raw []byte) (ret []types.Vehicle, err error) {
	err = json.Unmarshal(raw, &ret)
	return ret, err
}

func ProcessVehicles(cars []types.Vehicle) (vehicleOwnerIncomeInfo *types.VehicleOwnerIncomeInfo) {
	vehicleOwnerIncomeInfo = new(types.VehicleOwnerIncomeInfo)

	for k ,_:= range cars {
		vehicleOwnerIncomeInfo.Vehicles = append(vehicleOwnerIncomeInfo.Vehicles , &cars[k])
	}
	return vehicleOwnerIncomeInfo
}
/*
func getCarOwners() map[types.ID]*types.VehicleOwnerIncomeInfo{
	var carOwners = make(map[types.ID]*types.VehicleOwnerIncomeInfo)
	var filename = "/home/waris/car.json"
	bytes, err := readJson(filename)
	if err != nil {
		panic(err)
	}
	var leasedCars []types.Vehicle
	leasedCars, err = unmarshalLeasedCars(bytes)
	var carMap map[types.ID][]*types.Vehicle
	carMap = processLeasedCars(leasedCars)

	for carOwnerId, cars := range carMap {
		owner := types.VehicleOwnerIncomeInfo{}
		owner.OwnerId = carOwnerId
		owner.LeasedCars = make([]*types.Vehicle, 0)

		for _, car := range cars {
			owner.LeasedCars = append(owner.LeasedCars, car)
		}
		carOwners[carOwnerId] = &owner
	}
	return carOwners
	//printAny(carOwners)
}
*/
/*
func processLeasedCars(cars []types.Vehicle) (ret map[types.ID][]*types.Vehicle) {
	ret = make(map[types.ID][]*types.Vehicle)
	for k, car := range cars {
		if _, ok := ret[car.CarOwnerId]; !ok {
			ret[car.CarOwnerId] = make([]*types.Vehicle, 0)
		}

		ret[car.CarOwnerId] = append(ret[car.CarOwnerId], &cars[k])
	}
	return ret
}
*/

