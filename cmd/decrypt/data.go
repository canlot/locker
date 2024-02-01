/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package decrypt

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"main/internals"
)

var login string
var password string
var id string

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "decrypts data",
	Long: `Decrypts previously encrypted data with provided data id and login
Usage:
	decrypt data --id c711427a-0000-0000-8b93-54efa5d50310 --login user --password password`,
	Run: func(cmd *cobra.Command, args []string) {
		dataInfo, plainData, err := internals.DecryptData(id, login, password)
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
	dataCmd.Flags().StringVar(&password, "password", "", "Password for login")
	dataCmd.Flags().StringVar(&id, "id", "", "Data id")
	dataCmd.MarkFlagRequired("login")
	dataCmd.MarkFlagRequired("password")
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
