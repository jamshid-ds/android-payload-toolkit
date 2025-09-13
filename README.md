# Android Payload Toolkit 🚀

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-green)](https://github.com/yourusername/android-payload-toolkit)

A powerful Go toolkit for working with Android OTA payload.bin files. Extract partitions from OTA updates or create custom payload.bin files for Android devices.

## 🌟 Features

- **📦 Extract Mode**: Extract partition images from payload.bin files
- **🔨 Build Mode**: Create payload.bin from partition images  
- **⚡ Lightning Fast**: Parallel processing with customizable workers
- **📊 Progress Tracking**: Real-time progress bars for all operations
- **🔧 Flexible**: Support for unlimited number of partitions
- **🌐 Cross-Platform**: Works on Linux, macOS, and Windows

## 📋 Table of Contents

- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Usage](#-usage)
  - [Extract Mode](#extract-mode)
  - [Build Mode](#build-mode)
- [Examples](#-examples)
- [How It Works](#-how-it-works)
- [Building from Source](#-building-from-source)
- [Contributing](#-contributing)
- [License](#-license)

## 💾 Installation

### Prerequisites

- **XZ Utils** library (required for compression/decompression)

#### macOS
```bash
brew install xz
```

#### Linux (Ubuntu/Debian)
```bash
sudo apt-get install xz-utils liblzma-dev
```

#### Linux (CentOS/RHEL/Fedora)
```bash
sudo yum install xz-devel
# or
sudo dnf install xz-devel
```

#### Windows
Download and install from [XZ Utils official site](https://tukaani.org/xz/)

### Download Binary

Download the latest release for your platform from [Releases](https://github.com/yourusername/android-payload-toolkit/releases)

### Or Build from Source

```bash
git clone https://github.com/yourusername/android-payload-toolkit.git
cd android-payload-toolkit
go build -o android-payload-toolkit .
```

## 🚀 Quick Start

### Extract partitions from OTA update
```bash
# Extract all partitions
./android-payload-toolkit payload.bin -o extracted/

# Extract specific partitions
./android-payload-toolkit payload.bin -p boot,system,vendor -o extracted/
```

### Create payload.bin from images
```bash
./android-payload-toolkit build \
  -partitions boot:boot.img,system:system.img,vendor:vendor.img \
  -output custom_payload.bin
```

## 📖 Usage

### Extract Mode

Extract partition images from an existing payload.bin file.

```bash
./android-payload-toolkit [options] payload.bin
```

#### Options:
- `-o, --output` : Output directory for extracted images
- `-p, --partitions` : Comma-separated list of partitions to extract (extracts all if not specified)
- `-l, --list` : List partitions without extracting
- `-c, --concurrency` : Number of parallel workers (default: 4)

#### Examples:

```bash
# List all partitions in payload
./android-payload-toolkit payload.bin -l

# Extract all partitions
./android-payload-toolkit payload.bin -o output/

# Extract only boot and system
./android-payload-toolkit payload.bin -p boot,system -o output/

# Use 8 parallel workers for faster extraction
./android-payload-toolkit payload.bin -c 8 -o output/
```

### Build Mode

Create a new payload.bin file from partition images.

```bash
./android-payload-toolkit build [options]
```

#### Options:
- `-partitions` : **Required** - Comma-separated list of partition:image pairs
- `-output` : Output payload file (default: payload.bin)

#### Examples:

```bash
# Create payload with boot and system
./android-payload-toolkit build \
  -partitions boot:boot.img,system:system.img \
  -output payload.bin

# Create full OTA payload
./android-payload-toolkit build \
  -partitions boot:boot.img,system:system.img,vendor:vendor.img,product:product.img,\
system_ext:system_ext.img,odm:odm.img,recovery:recovery.img,dtbo:dtbo.img,\
vbmeta:vbmeta.img,vbmeta_system:vbmeta_system.img \
  -output full_ota.bin
```

## 💡 Examples

### Working with LineageOS

```bash
# Download LineageOS OTA
wget https://example.com/lineage-20.0-device.zip

# Extract payload.bin from ZIP
unzip -j lineage-20.0-device.zip payload.bin

# Extract all partitions
./android-payload-toolkit payload.bin -o lineage_parts/

# Modify boot.img (example: patch with Magisk)
# ...

# Rebuild payload with modified boot
./android-payload-toolkit build \
  -partitions boot:modified_boot.img,system:lineage_parts/system.img,\
vendor:lineage_parts/vendor.img \
  -output modified_payload.bin
```

### Working with Google Factory Images

```bash
# Download factory image
wget https://dl.google.com/dl/android/aosp/device-build-factory.zip

# Extract and find payload.bin
unzip device-build-factory.zip
cd device-build-factory
unzip image-device-build.zip payload.bin

# Extract partitions
../android-payload-toolkit payload.bin -o factory_parts/

# List extracted partitions
ls -lh factory_parts/
```

### Custom ROM Development

```bash
# Extract from multiple sources
./android-payload-toolkit source1.bin -o parts1/
./android-payload-toolkit source2.bin -o parts2/

# Combine partitions from different sources
./android-payload-toolkit build \
  -partitions boot:parts1/boot.img,\
system:custom_system.img,\
vendor:parts2/vendor.img,\
product:parts1/product.img \
  -output custom_rom.bin
```

## 🔄 How It Works

### Payload Structure

Android OTA payload.bin files follow a specific structure:

```
┌─────────────────────────┐
│    Magic ("CrAU")       │  4 bytes
├─────────────────────────┤
│    Version              │  8 bytes
├─────────────────────────┤
│    Manifest Size        │  8 bytes
├─────────────────────────┤
│  Metadata Signature Size│  4 bytes
├─────────────────────────┤
│    Manifest (protobuf)  │  Variable
├─────────────────────────┤
│    Metadata Signature   │  Variable
├─────────────────────────┤
│    Data Blobs           │  Variable
├─────────────────────────┤
│    Payload Signature    │  Variable
└─────────────────────────┘
```

### Workflow

#### Extract Mode:
1. **Parse Header**: Read magic, version, and metadata
2. **Load Manifest**: Deserialize protobuf manifest
3. **Process Partitions**: For each partition in manifest:
   - Read data blocks
   - Decompress if needed (REPLACE_XZ, REPLACE_BZ)
   - Write to output file
4. **Verify**: Check SHA256 hashes

#### Build Mode:
1. **Collect Images**: Read all partition image files
2. **Create Manifest**: Build protobuf with partition metadata
3. **Process Data**: For each partition:
   - Calculate SHA256 hash
   - Pad to block size (4096 bytes)
   - Store offset and size
4. **Write Payload**: Combine header + manifest + data blobs

## 🔨 Building from Source

### Requirements
- Go 1.18 or higher
- XZ development libraries
- Protocol Buffers compiler (optional, for modifying proto files)

### Build Steps

```bash
# Clone repository
git clone https://github.com/yourusername/android-payload-toolkit.git
cd android-payload-toolkit

# Install dependencies
go mod download

# Build for current platform
go build -o android-payload-toolkit .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o android-payload-toolkit-linux
GOOS=darwin GOARCH=amd64 go build -o android-payload-toolkit-macos
GOOS=windows GOARCH=amd64 go build -o android-payload-toolkit.exe
```

### macOS Build Note

On macOS, you may need to set CGO flags:

```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/xz/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/xz/lib"
go build -o android-payload-toolkit .
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Based on the Android Open Source Project update_engine
- Inspired by [payload-dumper-go](https://github.com/ssut/payload-dumper-go)
- Protocol definitions from ChromeOS Update Engine

## ⚠️ Disclaimer

This tool is for educational and development purposes only. Always respect device warranties and terms of service when modifying device firmware.

## 📞 Contact

- GitHub: [@yourusername](https://github.com/yourusername)
- Issues: [GitHub Issues](https://github.com/yourusername/android-payload-toolkit/issues)

---
Made with ❤️ for the Android development community