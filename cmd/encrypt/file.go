/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package encrypt

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internals"
)

var file string
var destination string

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Encrypt file",
	Long: `Encrypts file
Usage:
	locker encrypt file --source /var/file --destination /var/file.locked`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internals.EncryptFile(file, destination)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	fileCmd.Flags().StringVar(&file, "file", "", "source file")
	fileCmd.MarkFlagRequired("file")
	fileCmd.Flags().StringVar(&destination, "destination", "", "destination file")
	EncryptCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
