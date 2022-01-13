package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

const DefaultConfigTemplate = `###############################################################################
###                           Base Configuration                            ###
###############################################################################

# the in blockchain endpoint
in-endpoint = "http://localhost:8545"
# the out blockchain endpoint
out-endpoint = "http://localhost:8545"

# the bank contract address
bank = "0x4c2310DAdb5Be92a39336316f841e1944DA7bd60"

# the log fetching interval (milisec)
log-fetch-interval = 10000

###############################################################################
###                  Transaction Confirmer Configuration                    ###
###############################################################################
[confirmer]
# the total worker number
workers = 2
# the required minimum confirmation blocks
confirmation-blocks = 2
# the confirmation interval (milisec)
interval = 10
`

var configTemplate *template.Template

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initalize home directory",
	Long:  `initalize home directory creating the default config file`,
	Run: func(cmd *cobra.Command, args []string) {
		getConfig()
		os.RemoveAll(homeDir)
		err := os.Mkdir(homeDir, 0755)
		handleErr(err)
		path := homeDir + "/" + ConfigName + "." + ConfigType
		WriteConfigFile(path, nil)
		fmt.Printf("successfully initalized home directory: %s\n", path)
	},
}

func init() {
	var err error
	tmpl := template.New("appConfigFileTemplate")
	configTemplate, err = tmpl.Parse(DefaultConfigTemplate)
	handleErr(err)
	rootCmd.AddCommand(initCmd)
}

func WriteConfigFile(configFilePath string, config interface{}) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	mustWriteFile(configFilePath, buffer.Bytes(), 0644)
}

func mustWriteFile(filePath string, contents []byte, mode os.FileMode) {
	if err := ioutil.WriteFile(filePath, contents, mode); err != nil {
		fmt.Printf(fmt.Sprintf("failed to write file: %v", err) + "\n")
		os.Exit(1)
	}
}
