/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Displays data of one network interface",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		enteredIFace := args[0]

		iface, err := net.InterfaceByName(enteredIFace)

		if err != nil {
			fmt.Printf("[ERROR] getting network interface with name '%s': %s\n", enteredIFace, err.Error())
			os.Exit(1)
		}

		fmt.Printf("Name: %s\n", iface.Name)
		fmt.Printf("MTU (Maximum Transmission Unit): %d\n", iface.MTU)
		fmt.Printf("Hardware Address: %s\n", iface.HardwareAddr)
		fmt.Printf("Flags: %s\n", iface.Flags.String())
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
