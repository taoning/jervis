/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
    "fmt"
    "os"

	"github.com/spf13/cobra"
    "github.com/spf13/viper"
    "jervis/internal"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Edit the config file",
	Long: `Edit the config file`,
	Run: func(cmd *cobra.Command, args []string) {
        editor := os.Getenv("EDITOR")
        if editor == "" {
            editor = "vim"
        }
        configPath := viper.ConfigFileUsed()
        if configPath == "" {
            fmt.Println("Config file not found")
        }
        if err := internal.EditConfig(editor, configPath); err != nil {
            fmt.Println(err)
        }
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
