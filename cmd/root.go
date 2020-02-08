package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Sab94/go-watch-solidity/lib/generator"
	"github.com/spf13/cobra"
	"gopkg.in/fsnotify.v1"
)

var rootCmd = &cobra.Command{
	Use:   "go-watch-solidity",
	Short: "Go Watch Solidity",
	Long: `
   Go Watch Solidity is a watcher for a given solidity file.
   It generates abi, bin, and go bindings for the given solidity
   file on save.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		solCheck := exec.Command("solc", "--version")
		var out bytes.Buffer
		solCheck.Stdout = &out
		err := solCheck.Run()
		if err != nil {
			fmt.Println("solc not installed. Please install solc then try.")
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please specify a solidity file to watch")
			return
		}
		fileToWatch := args[0]
		if strings.HasSuffix(fileToWatch, ".sol") == false {
			fmt.Println("Please specify a solidity file to watch")
			return
		}
		if _, err := os.Stat(fileToWatch); err != nil {
			fmt.Println("Unable to find specified file")
			return
		}
		// Parse flags
		abi, _ := cmd.Flags().GetBool("abi")
		bin, _ := cmd.Flags().GetBool("bin")
		bindgo, _ := cmd.Flags().GetBool("bindgo")
		dest, _ := cmd.Flags().GetString("dest")

		fmt.Println("Watching : ", fileToWatch)

		SolidityWatcher(fileToWatch, abi, bin, bindgo, dest)
	},
}

// Execute executes the root command.
func Execute() error {
	rootCmd.PersistentFlags().BoolP("abi", "a", false, "Generate abi")
	rootCmd.PersistentFlags().BoolP("bin", "b", false, "Generate bin")
	rootCmd.PersistentFlags().BoolP("bindgo", "g", true, "Generate go binding")
	rootCmd.PersistentFlags().StringP("dest", "d", "", "Destination to generate")
	return rootCmd.Execute()
}

func SolidityWatcher(file string, abi bool, bin bool, bindgo bool, dest string) {
	// Generate once before starting Watcher
	err := generator.Generate(file, abi, bin, bindgo, dest)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println("Generated !!")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					err := generator.Generate(file, abi, bin, bindgo, dest)
					if err != nil {
						fmt.Println("error:", err)
						return
					}
					fmt.Println("Generated !!")
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(file)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
