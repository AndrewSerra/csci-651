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

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		interfaces, err := net.Interfaces()

		if err != nil {
			fmt.Printf("[ERROR] receiving network interfaces: %s\n", err)
			os.Exit(1)
		}

		for idx, iface := range interfaces {
			fmt.Printf("Interface Num %d\n", idx)
			fmt.Printf("\tName: %s\n", iface.Name)
			fmt.Printf("\tMTU (Maximum Transmission Unit): %d\n", iface.MTU)
			fmt.Printf("\tHardware Address: %s\n", iface.HardwareAddr)
			fmt.Printf("\tFlags: %s\n", iface.Flags.String())

			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}

			fmt.Printf("\tIP Addresses:\n")
			for _, addr := range addrs {
				fmt.Printf("\t\t%s\n", addr)
			}
			fmt.Println("-----------------------------")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
