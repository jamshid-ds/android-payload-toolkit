# Android Payload Toolkit - Detailed Usage Guide

## Table of Contents
1. [Understanding Payload Files](#understanding-payload-files)
2. [Extract Operations](#extract-operations)
3. [Build Operations](#build-operations)
4. [Advanced Usage](#advanced-usage)
5. [Troubleshooting](#troubleshooting)

## Understanding Payload Files

### What is payload.bin?
`payload.bin` is the core component of Android OTA (Over-The-Air) update packages. It contains:
- Partition images (boot, system, vendor, etc.)
- Metadata about each partition
- Checksums for verification
- Update instructions

### Where to find payload.bin?
1. **OTA Update ZIPs**: Inside any Android OTA update package
2. **Factory Images**: Google factory images contain payload.bin
3. **Custom ROMs**: LineageOS, PixelExperience, etc.

## Extract Operations

### Basic Extraction

Extract all partitions:
```bash
./android-payload-toolkit payload.bin -o output_folder/
```

### Selective Extraction

Extract only specific partitions:
```bash
# Extract only boot and system
./android-payload-toolkit payload.bin -p boot,system -o output/

# Extract critical partitions
./android-payload-toolkit payload.bin \
  -p boot,system,vendor,product,vbmeta \
  -o critical_parts/
```

### List Without Extracting

View available partitions:
```bash
./android-payload-toolkit payload.bin -l
```

### Performance Tuning

Adjust parallel workers for faster extraction:
```bash
# Use 8 workers (default is 4)
./android-payload-toolkit payload.bin -c 8 -o output/

# Use maximum available cores
./android-payload-toolkit payload.bin -c $(nproc) -o output/
```

## Build Operations

### Simple Build

Create payload from images:
```bash
./android-payload-toolkit build \
  -partitions boot:boot.img,system:system.img \
  -output custom.bin
```

### Full ROM Build

Build complete Android ROM payload:
```bash
./android-payload-toolkit build \
  -partitions boot:boot.img,\
init_boot:init_boot.img,\
vendor_boot:vendor_boot.img,\
dtbo:dtbo.img,\
vbmeta:vbmeta.img,\
vbmeta_system:vbmeta_system.img,\
vbmeta_vendor:vbmeta_vendor.img,\
recovery:recovery.img,\
system:system.img,\
system_ext:system_ext.img,\
product:product.img,\
vendor:vendor.img,\
vendor_dlkm:vendor_dlkm.img,\
odm:odm.img,\
odm_dlkm:odm_dlkm.img \
  -output full_rom.bin
```

### Dynamic Partition Building

Build from directory:
```bash
# Create partition list from directory
PARTITIONS=""
for img in images/*.img; do
  name=$(basename "$img" .img)
  PARTITIONS="${PARTITIONS}${name}:${img},"
done
PARTITIONS=${PARTITIONS%,}  # Remove trailing comma

# Build payload
./android-payload-toolkit build \
  -partitions "$PARTITIONS" \
  -output dynamic.bin
```

## Advanced Usage

### Workflow: Modify and Rebuild

1. **Extract original payload**
```bash
./android-payload-toolkit original.bin -o parts/
```

2. **Modify partitions**
```bash
# Example: Patch boot with Magisk
magiskboot unpack parts/boot.img
# ... modifications ...
magiskboot repack parts/boot.img
```

3. **Rebuild with modifications**
```bash
./android-payload-toolkit build \
  -partitions boot:parts/boot.img,\
system:parts/system.img,\
vendor:parts/vendor.img \
  -output modified.bin
```

### Merge Multiple Sources

Combine partitions from different ROMs:
```bash
# Extract from source ROMs
./android-payload-toolkit rom1.bin -o rom1_parts/
./android-payload-toolkit rom2.bin -o rom2_parts/

# Build hybrid ROM
./android-payload-toolkit build \
  -partitions boot:rom1_parts/boot.img,\
system:rom2_parts/system.img,\
vendor:rom1_parts/vendor.img,\
product:rom2_parts/product.img \
  -output hybrid.bin
```

### Verify Payload Integrity

Check payload without extracting:
```bash
# List partitions and sizes
./android-payload-toolkit payload.bin -l

# Extract and verify checksums
./android-payload-toolkit payload.bin -o verify/ 
sha256sum verify/*.img > checksums.txt
```

### Batch Processing

Process multiple payloads:
```bash
#!/bin/bash
for payload in *.bin; do
  echo "Processing $payload..."
  output_dir="${payload%.bin}_extracted"
  ./android-payload-toolkit "$payload" -o "$output_dir"
done
```

## Troubleshooting

### Common Issues

#### 1. "lzma.h not found" error
**Solution**: Install XZ development libraries
```bash
# macOS
brew install xz

# Ubuntu/Debian
sudo apt-get install liblzma-dev

# Fedora/RHEL
sudo dnf install xz-devel
```

#### 2. CGO compilation errors on macOS
**Solution**: Set environment variables
```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/xz/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/xz/lib"
make build
```

#### 3. "Verify failed" during extraction
**Cause**: Corrupted payload or incomplete download
**Solution**: Re-download the payload file

#### 4. Out of memory errors
**Solution**: Reduce concurrent workers
```bash
./android-payload-toolkit payload.bin -c 2 -o output/
```

### Performance Tips

1. **Use SSD**: Significantly faster than HDD
2. **Adjust workers**: More workers = faster, but uses more RAM
3. **Extract selectively**: Only extract needed partitions
4. **Clean temp files**: Remove old extractions to free space

### File Size Reference

Typical partition sizes:
- **boot.img**: 32-128 MB
- **system.img**: 2-4 GB
- **vendor.img**: 500 MB - 2 GB
- **product.img**: 200 MB - 1 GB
- **vbmeta.img**: 4-8 KB
- **dtbo.img**: 8-24 MB

Full payload sizes:
- **Pixel devices**: 2-3 GB
- **Samsung devices**: 4-6 GB
- **Custom ROMs**: 1.5-2.5 GB

## Best Practices

1. **Always backup original files** before modification
2. **Verify checksums** after extraction
3. **Test on virtual device** before flashing to physical device
4. **Keep partition order** consistent when rebuilding
5. **Document your modifications** for reproducibility

## Support

For issues and questions:
- GitHub Issues: [Report bugs](https://github.com/yourusername/android-payload-toolkit/issues)
- Discussions: [Ask questions](https://github.com/yourusername/android-payload-toolkit/discussions)