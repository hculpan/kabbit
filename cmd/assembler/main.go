/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/hculpan/kabbit/pkg/assembler"
	"github.com/hculpan/kabbit/pkg/executable"
	"github.com/spf13/cobra"
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "assembler <input file>",
	Short: "This assembler builds programs for the Kabbit VM",
	Long:  `This assembler builds programs for the Kabbit VM`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing required parameter: input file")
		}

		inputFile := args[0]
		outputFile, _ := cmd.Flags().GetString("output")
		debug, _ := cmd.Flags().GetBool("debug")
		if debug {
			fmt.Println("DEBUG output enabled")
		}

		if _, err := os.Stat(inputFile); err == nil {
			if len(outputFile) == 0 {
				ext := path.Ext(inputFile)
				outputFile = inputFile[0:len(inputFile)-len(ext)] + ".kbx"
			}
		} else {
			return errors.New(fmt.Sprintf("file '%s' not found", inputFile))
		}

		a := assembler.NewAssembler(debug)
		assembledCode, err := a.AssembleFromFile(inputFile)
		if err != nil {
			return err
		}

		header := assembledCode.NewFileHeader()
		ex := executable.NewExecutableFile(outputFile, header, assembledCode.Code, assembledCode.Data)
		fmt.Printf("Writing to output file %s\n", outputFile)
		ex.SaveFile()

		return nil
	},
	SilenceUsage: true,
}

func Execute() {
	fmt.Println("Kabbit Assembler v0.1.0")
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("output", "o", "", "Output file")
	rootCmd.Flags().BoolP("debug", "d", false, "Generate extra debug info")
}
