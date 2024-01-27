/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Adds new login, username can be empty",
	Long: `Adds new login, that login will encrypt private key
	Example:
		locker add login user
		locker add login
		username can stay empty.`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login called")
	},
}

func init() {
	AddCmd.AddCommand(loginCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
