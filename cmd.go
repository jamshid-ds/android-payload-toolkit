package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func buildPayload() {
	var (
		output     string
		partitions string
		inputDir   string
	)

	buildCmd := flag.NewFlagSet("build", flag.ExitOnError)
	buildCmd.StringVar(&output, "output", "payload.bin", "Output payload file")
	buildCmd.StringVar(&partitions, "partitions", "", "Comma-separated list of partition:image pairs")
	buildCmd.StringVar(&inputDir, "input", "", "Input directory containing partition images (auto-detect mode)")

	if len(os.Args) < 3 {
		printBuildUsage()
		os.Exit(1)
	}

	buildCmd.Parse(os.Args[2:])

	builder := NewPayloadBuilder(output)

	// Auto-detect mode: scan directory for .img files
	if inputDir != "" {
		if err := buildFromDirectory(builder, inputDir); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	} else if partitions != "" {
		// Manual mode: use specified partitions
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
	} else {
		fmt.Println("Error: Must specify either -input directory or -partitions list")
		printBuildUsage()
		os.Exit(1)
	}

	// Build payload
	if err := builder.Build(); err != nil {
		fmt.Printf("Error building payload: %v\n", err)
		os.Exit(1)
	}
}

func buildFromDirectory(builder *PayloadBuilder, dir string) error {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("directory not found: %s", dir)
	}

	// Read directory contents
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}

	// Count .img files
	imgCount := 0
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".img") {
			imgCount++
		}
	}

	if imgCount == 0 {
		return fmt.Errorf("no .img files found in directory: %s", dir)
	}

	fmt.Printf("Found %d partition images in %s\n", imgCount, dir)

	// Add each .img file as a partition
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".img") {
			// Remove .img extension to get partition name
			partName := strings.TrimSuffix(file.Name(), ".img")
			imagePath := filepath.Join(dir, file.Name())
			
			fmt.Printf("  Adding %s from %s\n", partName, file.Name())
			if err := builder.AddPartition(partName, imagePath); err != nil {
				return fmt.Errorf("failed to add partition %s: %v", partName, err)
			}
		}
	}

	return nil
}

func printBuildUsage() {
	fmt.Println("Usage: android-payload-toolkit build [options]")
	fmt.Println("\nOptions:")
	fmt.Println("  -output string")
	fmt.Println("        Output payload file (default \"payload.bin\")")
	fmt.Println("  -input string")
	fmt.Println("        Input directory containing .img files (auto-detect mode)")
	fmt.Println("  -partitions string")
	fmt.Println("        Comma-separated list of partition:image pairs (manual mode)")
	fmt.Println("\nExamples:")
	fmt.Println("  # Auto-detect all .img files in directory")
	fmt.Println("  android-payload-toolkit build -input extracted/ -output payload.bin")
	fmt.Println("")
	fmt.Println("  # Manual specification")
	fmt.Println("  android-payload-toolkit build -partitions boot:boot.img,system:system.img -output payload.bin")
	fmt.Println("")
	fmt.Println("  # With full paths")
	fmt.Println("  android-payload-toolkit build -partitions boot:/path/to/boot.img,system:/path/to/system.img")
}