#!/bin/bash
# Build custom payload.bin from partition images

# Check if required images exist
REQUIRED_IMAGES="boot.img system.img vendor.img"
MISSING=""

for img in $REQUIRED_IMAGES; do
    if [ ! -f "$img" ]; then
        MISSING="$MISSING $img"
    fi
done

if [ -n "$MISSING" ]; then
    echo "Error: Missing required images:$MISSING"
    echo "Please ensure all partition images are present"
    exit 1
fi

# Build minimal payload
echo "Building minimal payload with boot, system, and vendor..."
./android-payload-toolkit build \
    -partitions boot:boot.img,system:system.img,vendor:vendor.img \
    -output minimal_payload.bin

echo "Build complete: minimal_payload.bin"
ls -lh minimal_payload.bin