package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var disassemble bool

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "kabv <input file>",
	Short: "Executes programs in the Kabbit VM",
	Long:  `Executes programs in the Kabbit VM`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing required parameter: input file")
		}

		disassemble, _ = cmd.Flags().GetBool("disassemble")
		trace, _ := cmd.Flags().GetBool("trace")

		input := args[0]
		return ExecuteFile(input, disassemble, trace)
	},
	SilenceUsage: true,
}

func Execute() {
	fmt.Println("Kabbit Virtual Machine v0.1.0")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.Flags().StringP("output", "o", "", "Output file")
	rootCmd.Flags().BoolP("disassemble", "d", disassemble, "Output disassembly")
	rootCmd.Flags().BoolP("trace", "t", false, "Output trace information")
}
