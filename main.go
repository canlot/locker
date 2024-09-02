/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"main/cmd"
	"main/internals"
)

func main() {
	internals.CreateDatabaseIfNotExists()
	defer internals.Database.Close()
	cmd.Execute()
}
