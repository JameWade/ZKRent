package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ChengtayChain/ChengtayChain/crypto/ed25519"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

// sendtxCmd represents the sendtx command
var GenKeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate publickey and privatekey",
	Long:  `Generate publickey and privatekey`,
	RunE:  GenKey,
}

func GenKey(cmd *cobra.Command, args []string) error {
	var err error
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	var flagKeyDir string
	flag.StringVar(&flagKeyDir, "k", workingDir, "the directory which contains 'public.key.ed25519.json' file.")

	fmt.Println("Generating private key...")
	fmt.Println("Hint: if the program got stuck for more than a few seconds, you need to check the system's entropy source.")

	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()

	fmt.Println()

	{
		bytes, err := json.Marshal(privKey)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile("private.key.ed25519.json", bytes, 0644)
		if err != nil {
			panic(err)
		}
	}
	{
		bytes, err := json.Marshal(pubKey)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile("public.key.ed25519.json", bytes, 0644)
		if err != nil {
			return err
		}
	}

	return nil
}
