package show

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"main/internals"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show versions from software and database",
	Long:  "Show versions from software and database, also compatible database versions",
	Run: func(cmd *cobra.Command, args []string) {
		versions, err := internals.GetAllVersions()
		if err != nil {
			fmt.Println(err)
			return
		}
		tbl := table.New("Version", "Name", "Description")
		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		tbl.WithHeaderFormatter(headerFmt)
		for _, version := range versions {
			tbl.AddRow(version.Version, version.Name, version.Description)
		}
		tbl.Print()
	},
}

func init() {
	ShowCmd.AddCommand(versionCmd)
}
