/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package decrypt

import (
	"github.com/spf13/cobra"
)

var login string

// decryptCmd represents the decrypt command
var DecryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypts data or file",
	Long:  `Decrypt command decrypts previously encrypted file or data with one of existing logins`,
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
