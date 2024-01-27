/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"main/cmd"
	"main/internals"
)

func main() {
	cmd.Execute()
	internals.Database.Close()
}
