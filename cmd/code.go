/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jervis/internal"
)

// codeCmd represents the code command
var codeCmd = &cobra.Command{
	Use:   "code",
	Short: "Code clinic",
	Long: `Directly pass in code to have it diagnosed`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("readlines", readLines)
		viper.Set("format", format)
		viper.Set("sessionName", sessionName)
		viper.Set("newSession", newSession)
		internal.DoCode()
	},
}

func init() {
	rootCmd.AddCommand(codeCmd)
}
