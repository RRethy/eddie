package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/RRethy/eddie/internal/cmd/batch"
)

var (
	batchFile string
	batchJSON string
	batchOps  []string
)

var batchCmd = &cobra.Command{
	Use:   "batch",
	Short: "Execute multiple eddie operations in sequence from JSON input",
	Long: `Execute multiple eddie operations in sequence from JSON input.

Supports multiple input methods:
- From stdin: echo '{"operations":[...]}' | eddie batch
- From file: eddie batch --file operations.json  
- From JSON string: eddie batch --json '{"operations":[...]}'
- From operation flags: eddie batch --op view,file.txt --op str_replace,file.txt,old,new

Always continues execution on errors. Returns JSON output with success/error status for each operation.`,
	Run: func(cmd *cobra.Command, args []string) {
		var req *batch.BatchRequest
		var err error

		inputCount := 0
		if batchFile != "" {
			inputCount++
		}
		if batchJSON != "" {
			inputCount++
		}
		if len(batchOps) > 0 {
			inputCount++
		}
		if inputCount == 0 {
			inputCount = 1
		}

		if inputCount > 1 {
			fmt.Fprintln(os.Stderr, "Error: only one input method allowed")
			os.Exit(1)
		}

		switch {
		case batchFile != "":
			req, err = batch.ParseFromFile(batchFile)
		case batchJSON != "":
			req, err = batch.ParseFromJSON(batchJSON)
		case len(batchOps) > 0:
			req, err = batch.ParseFromOps(batchOps)
		default:
			req, err = batch.ParseFromStdin()
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing input: %v\n", err)
			os.Exit(1)
		}

		processor := batch.NewProcessor(os.Stdout)
		resp, err := processor.ProcessBatch(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing batch: %v\n", err)
			os.Exit(1)
		}

		output, err := json.Marshal(resp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling response: %v\n", err)
			os.Exit(1)
		}

		fmt.Print(string(output))
	},
}

func init() {
	rootCmd.AddCommand(batchCmd)
	batchCmd.Flags().StringVar(&batchFile, "file", "", "Read operations from JSON file")
	batchCmd.Flags().StringVar(&batchJSON, "json", "", "Operations as JSON string")
	batchCmd.Flags().StringArrayVar(&batchOps, "op", []string{}, "Individual operation (repeatable): type,arg1,arg2,...")
}
