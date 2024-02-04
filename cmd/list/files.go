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

// filesCmd represents the files command
var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		keys, hashes, fileInfo, err := internals.ListAllFiles()
		if err != nil {
			fmt.Println("Error occured")
			return
		} else {
			tbl := table.New("Key", "Hash", "Path", "Creation time")
			headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
			tbl.WithHeaderFormatter(headerFmt)
			for i := 0; i < len(keys); i++ {
				tbl.AddRow(keys[i], hashes[i], fileInfo[i].Path, fileInfo[i].CreateTime)
			}
			tbl.Print()
		}
	},
}

func init() {
	ListCmd.AddCommand(filesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// filesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// filesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
