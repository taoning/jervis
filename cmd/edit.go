/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jervis/internal"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit text",
	Long:  `Directly pass in text to have it edited.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("newSession", newSession)
        if sessionName != "" {
            viper.Set("editor.fileName", sessionName)
        }
		viper.Set("format", format)
		viper.Set("readlines", readLines)
		internal.DoEdit()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
