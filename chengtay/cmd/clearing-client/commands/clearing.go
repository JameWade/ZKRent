
package commands

import (
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/chengtay/cmd/clearing-client/internal"
	"github.com/ChengtayChain/ChengtayChain/chengtay/types"
	"github.com/shopspring/decimal"
	"strconv"

	"github.com/spf13/cobra"
)

// ClearingCmd represents the Clearing command
var ClearingCmd = &cobra.Command{
	Use:   "clearing",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: clearing,
}

func clearing(cmd *cobra.Command, args []string) error {
	VehicleOwnerId, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	_, totalPrice, chengtayPrice, _ := Compute(types.ID(VehicleOwnerId))
	//fmt.Println(vehicleOwnerIncomeInfo)
	fmt.Println("The VehicleOwner gets money:"+totalPrice)
	fmt.Println("The chengtay gets money:"+chengtayPrice)
	return nil
}
func Compute(vehicleOwnerId types.ID) (vehicleOwnerIncomeInfo *types.VehicleOwnerIncomeInfo, totalPrices string, chengtayPrices string, err error) {
	//50$ one hour
	hourPrice, err := decimal.NewFromString("50")
	chengtayPrice, _ := decimal.NewFromString("0")
	totalPrice, err := decimal.NewFromString("0")
	vehicleOwnerIncomeInfo = internal.GetVehicleOwner(vehicleOwnerId)

	if err != nil {
		panic(err)
	}
	//compute
	vehicles := vehicleOwnerIncomeInfo.Vehicles
	for _, vehicle := range vehicles {
		usedTime, _ := decimal.NewFromString(vehicle.UsedTime)
		fmt.Println("usedTime", usedTime)
		totalPrice = totalPrice.Add(hourPrice.Mul(usedTime))
	}
	//chengTay gets a proportional commission like 25%
	chengtayPrice = totalPrice.Mul(decimal.NewFromFloat32(0.25))
	////
	fmt.Println("totalPrice", totalPrice)
	vehicleOwnerIncomeInfo.Income = totalPrice.Sub(chengtayPrice).String()

	return vehicleOwnerIncomeInfo, totalPrice.String(), chengtayPrice.String(), nil
}
