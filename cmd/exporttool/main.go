package main

import (
	"fmt"
	"os"

	"github.com/temporalio/cloud-samples-go/export"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Example usage: exporttool /path/to/export/file")
		os.Exit(1)
	}

	filename := os.Args[1]

	bytes, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("error reading file: %v\r\n", err)
		os.Exit(1)
	}

	executions, err := export.DeserializeExportedExecutions(bytes)
	if err != nil {
		fmt.Printf("error extracting executions: %v\r\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully deserialized %d workflows and %d activities \r\n", len(executions.Workflows), len(executions.Activities))

	for _, workflow := range executions.Workflows {
		info, err := export.GetExportedWorkflowInformation(workflow)
		if err != nil {
			fmt.Printf("error extracting workflow information: %v\r\n", err)
			os.Exit(1)
		}

		fmt.Println(info)

		fmt.Println(export.FormatWorkflow(workflow))
		fmt.Println("----------------------------------------------------------")
		fmt.Println()
	}

	for _, activity := range executions.Activities {
		info, err := export.GetExportedActivityInformation(activity)
		if err != nil {
			fmt.Printf("error extracting activity information: %v\r\n", err)
			os.Exit(1)
		}

		fmt.Println(info)

		fmt.Println(export.FormatActivity(activity))
		fmt.Println("----------------------------------------------------------")
		fmt.Println()
	}
}
