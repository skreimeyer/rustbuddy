// Package cmd contains all core components of rustbuddy.
/*
Copyright © 2019 Samuel Kreimeyer

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rustbuddy [subcommand] [flags] [args...]",
	Short: "Rustbuddy is your friend. He fights for rustaceans",
	Long: `
  =========================================
  ___         _   _            _    _      
  | _ \_  _ __| |_| |__ _  _ __| |__| |_  _ 
  |   / || (_-<  _| '_ \ || / _' / _' | || |
  |_|_\\_,_/__/\__|_.__/\_,_\__,_\__,_|\_, |
                                       |__/ 
  =========================================
  __       __
  / <'     '> \
 (  / @   @ \  )
  \(_ _\_/_ _)/
(\ '-/     \-' /)
 "===\     /==="
  .==')___('==.    
 ' .='     '=.
 Rustbuddy is a small suite of tools to assist with rust development. Each 
 automates a small, but common task by code generation or static analysis. All
 generated code self-identifies with comments and aims for common-sense default
 behavior.
  
Happy oxidation!
`}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rustbuddy.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home := "$HOME"

		// Search config in home directory with name ".rustbuddy" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".rustbuddy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
