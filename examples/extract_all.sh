#!/bin/bash
# Extract all partitions from payload.bin

# Check if payload.bin exists
if [ ! -f "payload.bin" ]; then
    echo "Error: payload.bin not found!"
    echo "Please download an OTA update file and extract payload.bin"
    exit 1
fi

# Create output directory with timestamp
OUTPUT_DIR="extracted_$(date +%Y%m%d_%H%M%S)"

echo "Extracting all partitions to $OUTPUT_DIR"
./android-payload-toolkit payload.bin -o "$OUTPUT_DIR"

echo "Extraction complete!"
echo "Extracted files:"
ls -lh "$OUTPUT_DIR/"