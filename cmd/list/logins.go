/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internals"
)

// loginsCmd represents the logins command
var loginsCmd = &cobra.Command{
	Use:   "logins",
	Short: "Lists all logins",
	Long:  `Lists all logins with the login name and create date and time.`,
	Run: func(cmd *cobra.Command, args []string) {
		output, err := internals.ListAllUsers()
		if err != nil {
			fmt.Println("Error occured")
			return
		} else {
			fmt.Println(output)
		}

	},
}

func init() {
	ListCmd.AddCommand(loginsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
