/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"jervis/internal"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chatchat",
	Long:  `Starting a conversation.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("readlines", readLines)
		viper.Set("format", format)
		if sessionName != "" {
			viper.Set("chat.fileName", sessionName)
		}
		viper.Set("newSession", newSession)
		internal.DoChat()
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)
}
