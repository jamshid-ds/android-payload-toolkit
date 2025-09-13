package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func buildPayload() {
	var (
		output     string
		partitions string
	)

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&output, "output", "payload.bin", "Output payload file")
	buildCmd.StringVar(&partitions, "partitions", "", "Comma-separated list of partition:image pairs")

	if len(os.Args) < 3 {
		printBuildUsage()
		os.Exit(1)
	}

	buildCmd.Parse(os.Args[2:])

	if partitions == "" {
		fmt.Println("Error: No partitions specified")
		printBuildUsage()
		os.Exit(1)
	}

	builder := NewPayloadBuilder(output)

	// Parse partition list
	partList := strings.Split(partitions, ",")
	for _, part := range partList {
		parts := strings.Split(part, ":")
		if len(parts) != 2 {
			fmt.Printf("Error: Invalid partition format: %s\n", part)
			fmt.Println("Expected format: partition_name:image_file")
			os.Exit(1)
		}

		partName := strings.TrimSpace(parts[0])
		imagePath := strings.TrimSpace(parts[1])

		if err := builder.AddPartition(partName, imagePath); err != nil {
			fmt.Printf("Error adding partition %s: %v\n", partName, err)
			os.Exit(1)
		}
	}

	// Build payload
	if err := builder.Build(); err != nil {
		fmt.Printf("Error building payload: %v\n", err)
		os.Exit(1)
	}
}

func printBuildUsage() {
	fmt.Println("Usage: android-payload-toolkit build [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -output string")
	fmt.Println("        Output payload file (default \"payload.bin\")")
	fmt.Println("  -partitions string")
	fmt.Println("        Comma-separated list of partition:image pairs (required)")
	fmt.Println("\nExample:")
	fmt.Println("  android-payload-toolkit build -partitions boot:boot.img,system:system.img -output payload.bin")
}