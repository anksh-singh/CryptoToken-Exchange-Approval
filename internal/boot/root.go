package boot

import (
	"bridge-allowance/internal/adapters/bridge"
	"bridge-allowance/internal/adapters/cosmos"
	"bridge-allowance/internal/adapters/evm"
	"bridge-allowance/internal/adapters/nonevm"
	"bridge-allowance/web"
	"log"

	"github.com/spf13/cobra"
)

func init() {
	initAdapters()
}

var rootCmd = &cobra.Command{
	Use:   "gateway",
	Short: "bridge allowance",
	Long:  ` bridge allowance: Routes the http requests to multiple adapters`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf(err.Error())
	}
}
func initAdapters() {
	rootCmd.AddCommand(web.WebCmd)
	rootCmd.AddCommand(nonevm.SolanaCmd)
	rootCmd.AddCommand(evm.EvmCmd)
	rootCmd.AddCommand(cosmos.CosmosCmd)
	rootCmd.AddCommand(bridge.BridgeCmd)
}
