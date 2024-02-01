/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package add

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"main/internals"
	"syscall"
)

var loginFlag string

// var passwordFlag string
var loginNewFlag string

//var passwordNewFlag string

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Adds new login, username can be empty",
	Long: `Adds new login, that login will encrypt private key
	Example:
		locker add login --newlogin newuser 
		locker add login --login newuser --newlogin newuser2
		`,

	Run: func(cmd *cobra.Command, args []string) {
		empty, err := internals.IsDatabaseEmpty()
		if err != nil {
			fmt.Println("Internal error occur: " + err.Error())
			return
		}
		if (loginFlag == "") && !empty {
			fmt.Println("Database is not empty, existing login and password must be provided")
			return
		}
		if empty {
			fmt.Println("Password for login: ")
			bytePw, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			err = internals.CreateFirstLoginWithRSAKeys(loginNewFlag, string(bytePw))
			if err != nil {
				fmt.Println("Error occured at creation of first login: " + err.Error())
				return
			} else {
				fmt.Println("First login successful added")
				return
			}
		} else {
			fmt.Println("Password for existing login: ")
			existingPwd, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("Password for new login: ")
			newPwd, err := terminal.ReadPassword(int(syscall.Stdin))
			if err != nil {
				fmt.Println(err.Error())
			}
			err = internals.CreateLoginWithExistingRSAKeys(loginFlag, string(existingPwd), loginNewFlag, string(newPwd))
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
	//loginCmd.Flags().StringVar(&passwordFlag, "password", "", "Password for login")
	loginCmd.Flags().StringVar(&loginNewFlag, "newlogin", "", "New login name")
	//loginCmd.Flags().StringVar(&passwordNewFlag, "newpassword", "", "Password for new login")
	loginCmd.MarkFlagRequired("newlogin")
	//loginCmd.MarkFlagRequired("newpassword")
	AddCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
