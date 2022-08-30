/*
Copyright © 2022 NAME HERE BRADLEY.SOPER@CNVRG.IO

*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Will load a push docker images to the registry specified",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load command called")
		dockerPull()
	},
}

func init() {
	imagesCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dockerPull() {
	cmdStr := "echo hello"
	out, _ := exec.Command("/bin/sh", "-c", cmdStr).Output()
	fmt.Printf("%s", out)
}

/*
func readFile(f string) {
	content, err := ioutil.ReadFile(f)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(content))
}
*/
