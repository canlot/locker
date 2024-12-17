/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package decrypt

import (
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"main/internals"
	"syscall"
)

var source string
var destination string

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Decrypts file",
	Long: `Decrypts previously encrypted file with provided data id and login
Usage:
	locker decrypt file --source /var/file.lock
	locker decrypt file --source /var/file.lock --destination /opt`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		err = internals.DecryptFile(source, destination, login, string(password))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Data decrypted")
		}
	},
}

func init() {
	fileCmd.Flags().StringVarP(&source, "source", "s", "", "encrypted file")
	fileCmd.MarkFlagRequired("source")
	fileCmd.Flags().StringVarP(&destination, "destination", "d", "", "destination file")
	fileCmd.Flags().StringVarP(&login, "login", "l", "", "Login name")
	fileCmd.MarkFlagRequired("login")
	DecryptCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
