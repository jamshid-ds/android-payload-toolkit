# Android Payload Toolkit

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-green)](https://github.com/jamshid-ds/android-payload-toolkit)

Extract and build Android OTA payload.bin files with lightning-fast parallel processing.

## Features

- **Extract**: Unpack partition images from payload.bin files
- **Build**: Create custom payload.bin from partition images
- **Fast**: Parallel processing with customizable workers
- **Cross-platform**: Linux, macOS, and Windows support

## Installation

### Prerequisites
Install XZ Utils:
- **macOS**: `brew install xz`
- **Ubuntu/Debian**: `sudo apt-get install xz-utils liblzma-dev`
- **RHEL/Fedora**: `sudo dnf install xz-devel`
- **Windows**: [Download XZ Utils](https://tukaani.org/xz/)

### Download
Get the latest binary from [Releases](https://github.com/jamshid-ds/android-payload-toolkit/releases) or build from source:

```bash
git clone https://github.com/jamshid-ds/android-payload-toolkit.git
cd android-payload-toolkit
go build -o android-payload-toolkit .
```

## Quick Start

### Extract partitions
```bash
# Extract all partitions
./android-payload-toolkit payload.bin -o extracted/

# Extract specific partitions
./android-payload-toolkit payload.bin -p boot,system,vendor -o extracted/

# List partitions without extracting
./android-payload-toolkit payload.bin -l
```

### Build payload
```bash
# Auto-detect from directory
./android-payload-toolkit build -input extracted/ -output custom_payload.bin

# Specify partitions manually
./android-payload-toolkit build \
  -partitions boot:boot.img,system:system.img,vendor:vendor.img \
  -output custom_payload.bin
```

## Command Reference

### Extract Mode
```bash
./android-payload-toolkit [options] payload.bin
```

**Options:**
- `-o, --output`: Output directory
- `-p, --partitions`: Specific partitions to extract (comma-separated)
- `-l, --list`: List partitions only
- `-c, --concurrency`: Parallel workers (default: 4)

### Build Mode
```bash
./android-payload-toolkit build [options]
```

**Options:**
- `-input`: Directory with .img files (auto-detect mode)
- `-partitions`: Manual partition:image pairs
- `-output`: Output file (default: payload.bin)

> Use either `-input` OR `-partitions`, not both

## Examples

### Modify LineageOS boot image
```bash
# Extract partitions
./android-payload-toolkit lineage.bin -o parts/

# Patch boot.img with Magisk
# ...

# Rebuild with modified boot
./android-payload-toolkit build -input parts/ -output modified.bin
```

### Combine partitions from different ROMs
```bash
./android-payload-toolkit build \
  -partitions boot:rom1/boot.img,system:rom2/system.img,vendor:rom3/vendor.img \
  -output custom.bin
```

## Documentation

For detailed usage and examples, see [docs/USAGE.md](docs/USAGE.md)

## Building from Source

**Requirements:** Go 1.18+, XZ libraries

```bash
go mod download
go build -o android-payload-toolkit .
```

**macOS users** may need:
```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/xz/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/xz/lib"
```

## License

Apache License 2.0 - see [LICENSE](LICENSE)

## Acknowledgments

- Based on Android Open Source Project update_engine
- Inspired by [payload-dumper-go](https://github.com/ssut/payload-dumper-go)
