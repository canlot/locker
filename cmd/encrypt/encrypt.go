/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package encrypt

import (
	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var EncryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypts file or data",
	Long: `Encrypts file or data with public key without a use of a password, 
it creates a random password at the execute time and encrypts it with public key.`,
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
