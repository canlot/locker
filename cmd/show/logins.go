/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package show

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"main/internals"
)

// loginsCmd represents the logins command
var loginsCmd = &cobra.Command{
	Use:   "logins",
	Short: "Lists all logins",
	Long:  `Lists all logins with the login name and create date and time.`,
	Run: func(cmd *cobra.Command, args []string) {
		logins, err := internals.ListAllLogins()
		if err != nil {
			fmt.Println("Error occured")
			return
		}

		tbl := table.New("Login", "Creation time")
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt)
		for _, login := range logins {
			tbl.AddRow(login.Login, login.CreateTime)
		}
		tbl.Print()
	},
}

func init() {
	ShowCmd.AddCommand(loginsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
