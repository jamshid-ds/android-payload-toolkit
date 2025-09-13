package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"os"

	"google.golang.org/protobuf/proto"
	"github.com/yourusername/android-payload-toolkit/chromeos_update_engine"
)

type PayloadBuilder struct {
	outputFile   string
	partitions   map[string]string // partition name -> image file path
	manifest     *chromeos_update_engine.DeltaArchiveManifest
	writer       *os.File
	dataBlobs    [][]byte
	blobOffsets  []int64
	currentBlob  int64
	metadataHash hash.Hash
	payloadHash  hash.Hash
}

func NewPayloadBuilder(outputFile string) *PayloadBuilder {
	return &PayloadBuilder{
		outputFile:   outputFile,
		partitions:   make(map[string]string),
		dataBlobs:    make([][]byte, 0),
		blobOffsets:  make([]int64, 0),
		metadataHash: sha256.New(),
		payloadHash:  sha256.New(),
	}
}

func (pb *PayloadBuilder) AddPartition(name, imagePath string) error {
	if _, err := os.Stat(imagePath); err != nil {
		return fmt.Errorf("image file not found: %s", imagePath)
	}
	pb.partitions[name] = imagePath
	return nil
}

func (pb *PayloadBuilder) Build() error {
	// Create output file
	writer, err := os.Create(pb.outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer writer.Close()
	pb.writer = writer

	// Initialize manifest
	pb.manifest = &chromeos_update_engine.DeltaArchiveManifest{
		MinorVersion: proto.Uint32(0),
		Partitions: make([]*chromeos_update_engine.PartitionUpdate, 0),
	}

	// Process each partition
	for name, imagePath := range pb.partitions {
		if err := pb.processPartition(name, imagePath); err != nil {
			return fmt.Errorf("failed to process partition %s: %v", name, err)
		}
	}

	// Write payload header
	if err := pb.writeHeader(); err != nil {
		return err
	}

	// Write manifest
	if err := pb.writeManifest(); err != nil {
		return err
	}

	// Write data blobs
	if err := pb.writeDataBlobs(); err != nil {
		return err
	}

	fmt.Printf("Payload created successfully: %s\n", pb.outputFile)
	return nil
}

func (pb *PayloadBuilder) processPartition(name, imagePath string) error {
	fmt.Printf("Processing partition: %s from %s\n", name, imagePath)
	
	// Read image file
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return fmt.Errorf("failed to read image: %v", err)
	}

	// Pad data to block size for compatibility
	paddedData := data
	if len(data)%blockSize != 0 {
		padding := blockSize - (len(data) % blockSize)
		paddedData = append(data, make([]byte, padding)...)
	}
	
	// Calculate hash of padded data
	dataHash := sha256.Sum256(paddedData)

	// Create partition update
	partitionInfo := &chromeos_update_engine.PartitionInfo{
		Size: proto.Uint64(uint64(len(data))),
		Hash: dataHash[:],
	}

	// Calculate extent info - use actual data size, not block-aligned
	dataSize := uint64(len(data))
	extent := &chromeos_update_engine.Extent{
		StartBlock: proto.Uint64(0),
		NumBlocks:  proto.Uint64((dataSize + blockSize - 1) / blockSize),
	}

	// Create install operation
	operation := &chromeos_update_engine.InstallOperation{
		Type: chromeos_update_engine.InstallOperation_REPLACE.Enum(),
		DataOffset: proto.Uint64(uint64(pb.currentBlob)),
		DataLength: proto.Uint64(uint64(len(paddedData))),
		DstExtents: []*chromeos_update_engine.Extent{extent},
		DataSha256Hash: dataHash[:],
	}

	// Create partition update
	partition := &chromeos_update_engine.PartitionUpdate{
		PartitionName: proto.String(name),
		Operations: []*chromeos_update_engine.InstallOperation{operation},
		NewPartitionInfo: partitionInfo,
	}

	// Add to manifest
	pb.manifest.Partitions = append(pb.manifest.Partitions, partition)
	
	// Store blob data (padded for REPLACE operation)
	pb.dataBlobs = append(pb.dataBlobs, paddedData)
	pb.blobOffsets = append(pb.blobOffsets, pb.currentBlob)
	pb.currentBlob += int64(len(paddedData))

	return nil
}


func (pb *PayloadBuilder) writeHeader() error {
	// Write magic
	if _, err := pb.writer.Write([]byte(payloadHeaderMagic)); err != nil {
		return err
	}

	// Write version
	version := make([]byte, 8)
	binary.BigEndian.PutUint64(version, brilloMajorPayloadVersion)
	if _, err := pb.writer.Write(version); err != nil {
		return err
	}

	// Marshal manifest to get size
	manifestData, err := proto.Marshal(pb.manifest)
	if err != nil {
		return err
	}

	// Write manifest size
	manifestSize := make([]byte, 8)
	binary.BigEndian.PutUint64(manifestSize, uint64(len(manifestData)))
	if _, err := pb.writer.Write(manifestSize); err != nil {
		return err
	}

	// Write metadata signature size (0 for unsigned)
	sigSize := make([]byte, 4)
	binary.BigEndian.PutUint32(sigSize, 0)
	if _, err := pb.writer.Write(sigSize); err != nil {
		return err
	}

	return nil
}

func (pb *PayloadBuilder) writeManifest() error {
	manifestData, err := proto.Marshal(pb.manifest)
	if err != nil {
		return err
	}

	if _, err := pb.writer.Write(manifestData); err != nil {
		return err
	}

	return nil
}

func (pb *PayloadBuilder) writeDataBlobs() error {
	for i, blob := range pb.dataBlobs {
		fmt.Printf("Writing blob %d: %d bytes\n", i, len(blob))
		if _, err := pb.writer.Write(blob); err != nil {
			return err
		}
	}
	return nil
}