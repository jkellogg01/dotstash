/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dependCmd represents the depend command
var dependCmd = &cobra.Command{
	Use:     "depend [-r repository] path...",
	Aliases: []string{"add-config"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("depend called")
	},
}

func init() {
	rootCmd.AddCommand(dependCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dependCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dependCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
