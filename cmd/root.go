/*
Copyright © 2023 tnw

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jervis/internal/config"
)

var cfgFile string
var sessionName string
var newSession bool
var readLines bool
var format bool

var rootCmd = &cobra.Command{
	Use:   "jv",
	Short: "Chat completion",
	Long:  "Chat completion",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jervis.json)")
	rootCmd.PersistentFlags().StringVarP(&sessionName, "session", "s", "", "session name")
	rootCmd.PersistentFlags().BoolVarP(&newSession, "new", "n", false, "new session")
	rootCmd.PersistentFlags().BoolVarP(&readLines, "readlines", "r", false, "read lines")
	rootCmd.PersistentFlags().BoolVarP(&format, "format", "f", false, "format output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		config.SetDefaultConfig()
		// Search config in home directory with name ".jervis" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".jervis")
		viper.SafeWriteConfig()
	}

	viper.AutomaticEnv() // read in environment variables that match

	// viper.SetEnvPrefix("openai")
	viper.BindEnv("openai_api_key", "OPENAI_API_KEY")
	viper.BindEnv("anthropic_api_key", "ANTHROPIC_API_KEY")

	// viper.SetEnvPrefix("anthropic")
	// viper.BindEnv("anthropic_api_key")

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
}
