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

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit the config file",
	Long:  `Edit the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		editor := os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim"
		}
		configPath := viper.ConfigFileUsed()
		if configPath == "" {
			fmt.Println("Config file not found")
		}
		if err := config.EditConfig(editor, configPath); err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
