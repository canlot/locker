/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package encrypt

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internals"
)

var data string
var label string

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Encrypt data",
	Long: `Encrypts data
Usage:
	locker encrypt data --label datalabel --data secretdata`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internals.EncryptData(label, data)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	dataCmd.Flags().StringVar(&data, "data", "", "Plain data")
	dataCmd.MarkFlagRequired("data")
	dataCmd.Flags().StringVar(&label, "label", "", "Label for data")
	dataCmd.MarkFlagRequired("label")
	EncryptCmd.AddCommand(dataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
