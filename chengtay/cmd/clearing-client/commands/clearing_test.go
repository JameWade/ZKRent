package commands

import (
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/cmd/clearing-client/internal"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"io/ioutil"
	"reflect"
	"testing"
)
//func GenTestVehicleData() (TestFile string){
//	vehicle := []types.Vehicle{
//		{"1", "0","1"},
//		{"1", "1","2"},
//		{"1","2","3"},
//	}
//	leasedCarsBytes, err := json.Marshal(vehicle)
//	if err != nil{
//		panic(err)
//	}
//	homeDir, err := homedir.Dir()
//	TestFile = homeDir + string(os.PathSeparator) + "carinfo.json"
//	err = ioutil.WriteFile(TestFile,leasedCarsBytes,0666)
//	return TestFile
//}

func Equal(t *testing.T, exp, got interface{}) {
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("\033[31m\n exp: %v \n got: %v \033[39m\n", exp, got)
		t.FailNow()
	}
}

func TestCompute(t *testing.T) {
	carowner,totalPrice,chengtayPrice,_ := Compute("1")
	t.Log(carowner)
	Equal(t,totalPrice,"300")
	Equal(t,chengtayPrice,"75")
}

func TestProcessLeasedCars(t *testing.T)  {
	filename := internal.GenTestVehicleData()
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var leasedCars []types.Vehicle
	leasedCars, err = internal.UnmarshalVehicles(bytes)
	fmt.Printf("+%v",leasedCars[0])
	carOwner := internal.ProcessVehicles(leasedCars)
	t.Log(carOwner)
}
func TestApi(t *testing.T) {
	vehicleOwnerIncomeInfo := internal.GetVehicleOwner("1")
	t.Log(vehicleOwnerIncomeInfo)
	for _,car :=range vehicleOwnerIncomeInfo.Vehicles{
		t.Log(car)
	}
}


