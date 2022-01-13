package bridgecli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tak1827/evm-bridge/cli/log"
)

const (
	EnvPrefix       = "bridgecli"
	ConfigType      = "toml"
	ConfigName      = "config"
	DefaultHomeName = ".bridgecli"
)

var (
	cfgFile string
	homeDir string
	logger  = log.CLI("")
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bridgecli",
	Short: "The cli tool of bridging ERC20 and NFT",
	Long: `The comand line tool of server side implementaion of ERC20 and NFT bridging implementaion.
Fetch contract events periodically, then mint new asset to the other chain.`,
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()
		fmt.Println("use `-h` option to see help")
	},
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

	rootCmd.PersistentFlags().StringVar(&homeDir, "home", "", fmt.Sprintf("the home directory path (default is $HOME/%s/)", DefaultHomeName))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// viper.AddConfigPath(".")
		viper.SetConfigType(ConfigType)
		viper.SetConfigName(ConfigName)
	}

	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
}

func handleErr(err error) {
	if err != nil {
		logger.Error().Stack().Err(err).Msg("failed serving")
		logger.Fatal().Msg("stop serving")
	}
}

func getConfig() {
	if homeDir == "" {
		home, err := os.UserHomeDir()
		handleErr(err)
		homeDir = home + "/" + DefaultHomeName
	}
	viper.AddConfigPath(homeDir)

	if err := viper.ReadInConfig(); err == nil {
		// fmt.Printf("using config file: %s\n", viper.ConfigFileUsed())
	}
}

func getConfigString(key string, val *string) {
	if *val == "" {
		if *val = viper.GetString(key); *val == "" {
			logger.Fatal().Msgf("no `%s` setting", key)
		}
	}
	logger.Info().Msgf("%s: %s", key, *val)
}

func getConfigInt(key string, val *int) {
	if *val == 0 {
		if *val = viper.GetInt(key); *val == 0 {
			logger.Fatal().Msgf("no `%s` setting", key)
		}
	}
	logger.Info().Msgf("%s: %d", key, *val)
}
