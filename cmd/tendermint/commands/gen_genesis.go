package commands

import (
	"fmt"
	tmos "github.com/ChengtayChain/ChengtayChain/libs/os"
	tmrand "github.com/ChengtayChain/ChengtayChain/libs/rand"
	"github.com/ChengtayChain/ChengtayChain/privval"
	"github.com/ChengtayChain/ChengtayChain/types"
	tmtime "github.com/ChengtayChain/ChengtayChain/types/time"
	"github.com/spf13/cobra"
	"strconv"
)

var GenGenesisCmd = &cobra.Command{
	Use:   "gen_genesis [number of validators]",
	Short: "Generate a new genesis file and specified number of validators, one of which is the top validator with power 100 while the rest has power 50",
	Args:  cobra.ExactArgs(1),
	RunE:  genGenesis,
}

func genGenesis(cmd *cobra.Command, args []string) error {
	validatorNum, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	if validatorNum == 0 {
		err := fmt.Errorf("At least one validator is required")
		logger.Error(err.Error())
		return err
	}

	// load genesis file
	genFile := config.GenesisFile()
	if tmos.FileExists(genFile) {
		err = fmt.Errorf("Genesis file already exists")
		logger.Error(err.Error(), "path", genFile)
		return err
	}

	// create private keys for multiple validator
	validators := make([]types.GenesisValidator, 0)

	for i := 0; i < validatorNum; i++ {
		privValKeyFile := config.PrivValidatorKeyFile() + "." + fmt.Sprint(i) + ".json"
		privValStateFile := config.PrivValidatorStateFile() + "." + fmt.Sprint(i) + ".json"

		pv := privval.GenFilePV(privValKeyFile, privValStateFile)
		pubkey, err := pv.GetPubKey()
		if err != nil {
			return err
		}
		pv.Save()

		validator := types.GenesisValidator{
			Address: pubkey.Address(),
			PubKey:  pubkey,
			Power:   0,
			Name:    "Genesis validator " + fmt.Sprint(i),
		}

		if i == 0 {
			validator.Power = 100
		} else {
			validator.Power = 50
		}

		validators = append(validators, validator)

		logger.Info("Generated private validator", "keyFile", privValKeyFile,
			"stateFile", privValStateFile, "validator", i, "power", validator.Power)
	}

	genDoc := types.GenesisDoc{
		ChainID:         fmt.Sprintf("chengtay-chain-%v", tmrand.Str(6)),
		GenesisTime:     tmtime.Now(),
		ConsensusParams: types.DefaultConsensusParams(),
	}

	genDoc.Validators = validators

	if err := genDoc.SaveAs(genFile); err != nil {
		return err
	}
	logger.Info("Generated genesis file", "path", genFile)

	logger.Info("Please distribute these validator key files and validator state files for each validators. These files are secret.")

	return nil
}
