// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type Config struct {
	Alias   string `yaml:"alias"`
	Event   string `yaml:"event"`
	Channel string `yaml:"channel"`
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "create a config file in your home directory",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		home, err := homedir.Dir()
		if err != nil {
			fmt.Println("Unable to get Home Directory: ", err)
			os.Exit(1)
		}

		cfgf := filepath.Join(home, ".cda.yaml")
		_, err = os.Stat(cfgf)
		if err == nil {
			fmt.Println("Config file exists, please delete before creating again:", cfgf)
			return
		}
		config := []byte("Alias: replaceme\n")
		err = ioutil.WriteFile(cfgf, config, 0644)
		if err != nil {
			fmt.Println("Error creating config file: ", err)
		} else {
			fmt.Println("Config file created at ", cfgf)
			fmt.Println("Please edit this file with your microsoft alias")
		}

	},
}

func init() {
	RootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
