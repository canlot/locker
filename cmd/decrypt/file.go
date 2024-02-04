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
var loginFile string

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Password: ")
		password, err := term.ReadPassword(int(syscall.Stdin))
		err = internals.DecryptFile(source, destination, loginFile, string(password))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("Data decrypted")
		}
	},
}

func init() {
	fileCmd.Flags().StringVar(&source, "source", "", "source file")
	fileCmd.MarkFlagRequired("source")
	fileCmd.Flags().StringVar(&destination, "destination", "", "destination file")
	fileCmd.MarkFlagRequired("destination")
	fileCmd.Flags().StringVar(&loginFile, "login", "", "Login name")
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
