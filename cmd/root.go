/*
Copyright © 2020
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// welcomeMsg is a multiline string that represents the ASCII art of the sport info CLI logo.
var welcomeMsg = `
 ___  ____  _____  ____  ____    ____  _  _  ____  _____
/ __)(  _ \(  _  )(  _ \(_  _)  (_  _)( \( )( ___)(  _  )
\__ \ )___/ )(_)(  )   /  )(     _)(_  )  (  )__)  )(_)(
(___/(__)  (_____)(_)\_) (__)   (____)(_)\_)(__)  (_____)

`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sport-info",
	Short: " Sport info CLI",
	Long:  welcomeMsg + "The CLI wizard\n\n" + "Welcome to the Sport info CLI!",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
