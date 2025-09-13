# Installation Guide

## 🚀 Quick Install (Recommended)

### Download Pre-built Binary

1. Go to [Releases](https://github.com/jamshid-ds/android-payload-toolkit/releases)
2. Download the binary for your platform:
   - macOS M1/M2: `android-payload-toolkit-darwin-arm64`
   - macOS Intel: `android-payload-toolkit-darwin-amd64`
   - Linux: `android-payload-toolkit-linux-amd64`
   - Windows: `android-payload-toolkit-windows-amd64.exe`

3. Make it executable (macOS/Linux):
```bash
chmod +x android-payload-toolkit-*
```

4. Run directly:
```bash
./android-payload-toolkit-darwin-arm64 --help
```

## 🔨 Build from Source (Optional)

### Prerequisites

Install XZ library:

**macOS:**
```bash
brew install xz
```

**Linux:**
```bash
# Ubuntu/Debian
sudo apt-get install xz-utils liblzma-dev

# Fedora/RHEL
sudo dnf install xz-devel
```

### Build

```bash
# Clone
git clone https://github.com/jamshid-ds/android-payload-toolkit.git
cd android-payload-toolkit

# Build (macOS)
export CGO_CFLAGS="-I/opt/homebrew/opt/xz/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/xz/lib"
go build -o android-payload-toolkit .

# Build (Linux)
go build -o android-payload-toolkit .
```

## ✅ Verify Installation

```bash
./android-payload-toolkit --help
```

## 📝 Usage

### Extract
```bash
./android-payload-toolkit payload.bin -o extracted/
```

### Build from folder
```bash
./android-payload-toolkit build -input extracted/ -output new.bin
```

That's it! No `make` needed if you download the binary.