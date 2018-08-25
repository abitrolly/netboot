// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go.universe.tf/netboot/pixiecore"
	"go.universe.tf/netboot/pixiecore/cli"
	"go.universe.tf/netboot/third_party/ipxe"
)

func initConfig() {
	viper.SetEnvPrefix("pixiecore")
	viper.AutomaticEnv()
}

var rootCmd = &cobra.Command{
	Use:   "netboot.xyz",
	Short: "netboot.xyz PXE boot server",
	Run: func(cmd *cobra.Command, args []string) {
		kernel := "https://boot.netboot.xyz/ipxe/netboot.xyz.lkrn"
		fmt.Println(cli.StaticFromFlags(cmd, kernel, []string{}, "").Serve())
	},
}

func main() {
	cli.Ipxe[pixiecore.FirmwareX86PC] = ipxe.MustAsset("undionly.kpxe")
	cli.Ipxe[pixiecore.FirmwareEFI32] = ipxe.MustAsset("ipxe-i386.efi")
	cli.Ipxe[pixiecore.FirmwareEFI64] = ipxe.MustAsset("ipxe-x86_64.efi")
	cli.Ipxe[pixiecore.FirmwareEFIBC] = ipxe.MustAsset("ipxe-x86_64.efi")

	rootCmd.Flags().String("cmdline", "", "Kernel commandline arguments")
	rootCmd.Flags().String("bootmsg", "", "Message to print on machines before booting")
	
	rootCmd.Flags().BoolP("debug", "d", false, "Log more things that aren't directly related to booting a recognized client")
	rootCmd.Flags().BoolP("log-timestamps", "t", false, "Add a timestamp to each log line")
	rootCmd.Flags().StringP("listen-addr", "l", "0.0.0.0", "IPv4 address to listen on")
	rootCmd.Flags().IntP("port", "p", 80, "Port to listen on for HTTP")
	rootCmd.Flags().Int("status-port", 0, "HTTP port for status information (can be the same as --port)")
	
	// this is the default
	rootCmd.Flags().Bool("dhcp-no-bind", true, "Handle DHCP traffic without binding to the DHCP server port")
	rootCmd.Flags().String("ipxe-bios", "", "Path to an iPXE binary for BIOS/UNDI")
	rootCmd.Flags().String("ipxe-efi32", "", "Path to an iPXE binary for 32-bit UEFI")
	rootCmd.Flags().String("ipxe-efi64", "", "Path to an iPXE binary for 64-bit UEFI")

	// Development flags, hidden from normal use.
	rootCmd.Flags().String("ui-assets-dir", "", "UI assets directory (used for development)")
	rootCmd.Flags().MarkHidden("ui-assets-dir")

	cobra.OnInitialize(initConfig)	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	os.Exit(0)
}
