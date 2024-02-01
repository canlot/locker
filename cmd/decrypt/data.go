/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package decrypt

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"main/internals"
	"syscall"
)

var login string
var id string

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "decrypts data",
	Long: `Decrypts previously encrypted data with provided data id and login
Usage:
	decrypt data --id c711427a-0000-0000-8b93-54efa5d50310 --login user`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Password: ")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		dataInfo, plainData, err := internals.DecryptData(id, login, string(password))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			tbl := table.New("Label", "Decrypted data", "Creation time")
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			tbl.WithHeaderFormatter(headerFmt)
			tbl.AddRow(dataInfo.Label, plainData, dataInfo.CreateTime)
			tbl.Print()
		}

	},
}

func init() {
	dataCmd.Flags().StringVar(&login, "login", "", "Login name")
	dataCmd.Flags().StringVar(&id, "id", "", "Data id")
	dataCmd.MarkFlagRequired("login")
	dataCmd.MarkFlagRequired("id")
	DecryptCmd.AddCommand(dataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
