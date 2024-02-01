/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internals"
)

var loginFlag string
var passwordFlag string
var loginNewFlag string
var passwordNewFlag string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Adds new login, username can be empty",
	Long: `Adds new login, that login will encrypt private key
	Example:
		locker add login --newlogin newuser --newpassword password12345
		locker add login --login newuser --password password12345 --newlogin newuser2 --newpassword password98765
		`,

	Run: func(cmd *cobra.Command, args []string) {
		empty, err := internals.IsDatabaseEmpty()
		if err != nil {
			fmt.Println("Internal error occur: " + err.Error())
			return
		}
		if (loginFlag == "" || passwordFlag == "") && !empty {
			fmt.Println("Database is not empty, existing login and password must be provided")
			return
		}
		if empty {
			err = internals.CreateFirstLoginWithRSAKeys(loginNewFlag, passwordNewFlag)
			if err != nil {
				fmt.Println("Error occured at creation of first login: " + err.Error())
				return
			} else {
				fmt.Println("First login successful added")
				return
			}
		} else {
			err = internals.CreateLoginWithExistingRSAKeys(loginFlag, passwordFlag, loginNewFlag, passwordNewFlag)
			if err != nil {
				fmt.Println("Error occured at creation of login: " + err.Error())
				return
			} else {
				fmt.Println("Login succesfull added")
			}
		}
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginFlag, "login", "", "Login name")
	loginCmd.Flags().StringVar(&passwordFlag, "password", "", "Password for login")
	loginCmd.Flags().StringVar(&loginNewFlag, "newlogin", "", "New login name")
	loginCmd.Flags().StringVar(&passwordNewFlag, "newpassword", "", "Password for new login")
	loginCmd.MarkFlagRequired("newlogin")
	loginCmd.MarkFlagRequired("newpassword")
	AddCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
