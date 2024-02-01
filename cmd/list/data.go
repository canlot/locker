/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package list

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"main/internals"
)

// dataCmd represents the data command
var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		keys, dataInfo, err := internals.ListAllData()
		if err != nil {
			fmt.Println("Error occured")
			return
		} else {
			tbl := table.New("Key", "Label", "Creation time")
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			tbl.WithHeaderFormatter(headerFmt)
			for i := 0; i < len(keys); i++ {
				tbl.AddRow(keys[i], dataInfo[i].Label, dataInfo[i].CreateTime)
			}
			tbl.Print()
		}
	},
}

func init() {
	ListCmd.AddCommand(dataCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dataCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dataCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
