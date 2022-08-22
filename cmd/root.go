package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/littlehawk93/csc364asm/asm"
	"github.com/littlehawk93/ihex"
	"github.com/spf13/cobra"
)

const (
	instructionByteSize = 2
)

var inputFilePath string
var outputFilePath string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "asm364",
	Short: "CLI assembler for the Louisiana Tech University CSC 364 assembly language",
	Run:   runRootCommand,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init sets the command line parameters used by this program
func init() {

	log.SetFlags(log.Ldate | log.Ltime)

	rootCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "An input file to read assembly instructions from")
	rootCmd.Flags().StringVarP(&inputFilePath, "output", "o", "", "An output file to write binary HEX instructions to")
}

// runRootCommand the "main" method for running the root command
func runRootCommand(cmd *cobra.Command, args []string) {

	inputReader := os.Stdin

	// if an input file is provided, read assembly instructions from the file instead of stdin
	if inputFilePath != "" {
		f, err := os.Open(inputFilePath)

		if err != nil {
			log.Fatalf("Unable to open file '%s': %s\n", inputFilePath, err.Error())
		}

		defer f.Close()
		inputReader = f
	}

	outputWriter := os.Stdout

	// if an output file is provided, write assembly binary to the file instead of stdout
	if outputFilePath != "" {
		f, err := os.Create(outputFilePath)

		if err != nil {
			log.Fatalf("Unable to create file '%s': %s\n", outputFilePath, err.Error())
		}

		defer f.Close()
		outputWriter = f
	}

	writer, err := ihex.NewFileWriter(outputWriter, instructionByteSize)

	if err != nil {
		log.Fatalf("Error opening HEX writer: %s\n", err.Error())
	}

	defer writer.Close()

	parser := asm.NewParser(inputReader)

	buf := &bytes.Buffer{}

	// write all of the instruction binary data to a temporary buffer
	// if there are any syntax errors, the output writer is never written to
	for d, ok, err := parser.Next(); ok; d, ok, err = parser.Next() {
		if err != nil {
			log.Fatalln(err.Error())
		}

		if _, err = buf.Write(d); err != nil {
			log.Fatalf("Error writing buffer data: %s\n", err.Error())
		}
	}

	// copy buffer data into the output writer
	if _, err = writer.Write(buf.Bytes()); err != nil {
		log.Fatalf("Error writing assembly binary file: %s\n", err.Error())
	}
}
