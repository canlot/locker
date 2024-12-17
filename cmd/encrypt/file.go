/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package encrypt

import (
	"fmt"
	"github.com/spf13/cobra"
	"main/internals"
)

var source string
var destination string

// fileCmd represents the file command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Encrypt file",
	Long: `Encrypts file
Usage:
	locker encrypt file --file /var/file --destination /var/file.locked`,
	Run: func(cmd *cobra.Command, args []string) {
		err := internals.EncryptFile(source, destination)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	fileCmd.Flags().StringVarP(&source, "source", "s", "", "source file")
	fileCmd.MarkFlagRequired("source")
	fileCmd.Flags().StringVarP(&destination, "destination", "d", "", "destination file or path")
	EncryptCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
