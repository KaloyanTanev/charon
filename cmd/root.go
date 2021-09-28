/*
Copyright © 2021 Oisín Kyne <oisin@obol.tech>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var beaconNodes string
var peerNodes string
var quiet bool
var verbose bool
var debug bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "charon",
	Short: "Charon - The Ethereum SSV middleware client",
	Long: `Charon client(s) enable the division of Ethereum validator operation across a group of trusted parties using threshold cryptography.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) { 
		fmt.Println("No command specified, starting Charon as an SSV client")
	},
	PersistentPreRunE: persistentPreRunE,
}

// Pre-run hook executed by all commands and subcommands unless they declare their own
// Used to parse and validate global config typically (for now it sets log level)
func persistentPreRunE(cmd *cobra.Command, args []string) error {
	if cmd.Name() == "help" {
		// User just wants help
		return nil
	}

	if cmd.Name() == "version" {
		// User just wants the version
		return nil
	}

	// Disable service logging.
	// zerolog.SetGlobalLevel(zerolog.Disabled)

	// We bind viper here so that we bind to the correct command.
	quiet = viper.GetBool("quiet")
	verbose = viper.GetBool("verbose")
	debug = viper.GetBool("debug")
	// // Command-specific bindings.
	// switch fmt.Sprintf("%s/%s", cmd.Parent().Name(), cmd.Name()) {
	// case "account/create":
	// 	accountCreateBindings()
	// case "account/derive":
	// 	accountDeriveBindings()
	// case "account/import":
	// 	accountImportBindings()
	// case "attester/duties":
	// 	attesterDutiesBindings()
	// case "attester/inclusion":
	// 	attesterInclusionBindings()
	// case "block/info":
	// 	blockInfoBindings()
	// case "chain/time":
	// 	chainTimeBindings()
	// case "exit/verify":
	// 	exitVerifyBindings()
	// case "node/events":
	// 	nodeEventsBindings()
	// case "slot/time":
	// 	slotTimeBindings()
	// case "synccommittee/members":
	// 	synccommitteeMembersBindings()
	// case "validator/depositdata":
	// 	validatorDepositdataBindings()
	// case "validator/duties":
	// 	validatorDutiesBindings()
	// case "validator/exit":
	// 	validatorExitBindings()
	// case "validator/info":
	// 	validatorInfoBindings()
	// case "validator/keycheck":
	// 	validatorKeycheckBindings()
	// case "wallet/create":
	// 	walletCreateBindings()
	// case "wallet/import":
	// 	walletImportBindings()
	// case "wallet/sharedexport":
	// 	walletSharedExportBindings()
	// case "wallet/sharedimport":
	// 	walletSharedImportBindings()
	// }

	if quiet && verbose {
		fmt.Println("Cannot supply both quiet and verbose flags")
	}
	if quiet && debug {
		fmt.Println("Cannot supply both quiet and debug flags")
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.charon.yaml)")
	rootCmd.PersistentFlags().StringVar(&beaconNodes, "beacon-node", "http://localhost:5051", "URI for beacon node API")
	rootCmd.PersistentFlags().StringVar(&peerNodes, "peers", "http://localhost:9001,http://localhost:9002,http://localhost:9003", "URIs of peer charon clients")
	
	rootCmd.PersistentFlags().Bool("quiet", false, "do not generate any output")
	if err := viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet")); err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().Bool("verbose", false, "generate additional output where appropriate")
	if err := viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose")); err != nil {
		panic(err)
	}
	rootCmd.PersistentFlags().Bool("debug", false, "generate debug output")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".charon" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".charon")
	}

	viper.SetEnvPrefix("CHARON")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		// Don't report lack of config file...
		assert(strings.Contains(err.Error(), "Not Found"), "failed to read configuration")
	}
}
