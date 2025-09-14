# Android Payload Toolkit 🚀

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat\&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-green)](https://github.com/jamshid-ds/android-payload-toolkit)

Toolkit for working with **Android OTA payload.bin**: extract partitions or build custom payload files.

---

## 🌟 Features

* 📦 Extract partition images from payload.bin
* 🔨 Build payload.bin from .img files
* ⚡ Fast parallel processing
* 📊 Progress bars
* 🌐 Cross-platform

---

## 💾 Installation

### Dependencies

* **XZ Utils** required

  * macOS: `brew install xz`
  * Ubuntu/Debian: `sudo apt-get install xz-utils liblzma-dev`
  * Fedora: `sudo dnf install xz-devel`
  * Windows: [XZ Utils](https://tukaani.org/xz/)

### Build from Source

```bash
git clone https://github.com/jamshid-ds/android-payload-toolkit.git
cd android-payload-toolkit
go build -o android-payload-toolkit .
```

---

## 🚀 Quick Start

### Extract

```bash
# Extract all
./android-payload-toolkit payload.bin -o extracted/

# Extract selected
./android-payload-toolkit payload.bin -p boot,system -o extracted/
```

### Build

```bash
# Auto-detect from folder
./android-payload-toolkit build -input extracted/ -output custom.bin

# Manual partitions
./android-payload-toolkit build \
  -partitions boot:boot.img,system:system.img \
  -output custom.bin
```

---

## 🔨 Build Notes

* Go 1.18+ required
* For macOS:

```bash
export CGO_CFLAGS="-I/opt/homebrew/opt/xz/include"
export CGO_LDFLAGS="-L/opt/homebrew/opt/xz/lib"
go build -o android-payload-toolkit .
```

---

## 📄 License

Apache 2.0 – see [LICENSE](LICENSE)

## 📞 Contact

* GitHub: [@jamshid-ds](https://github.com/jamshid-ds)
* Issues: [Report here](https://github.com/jamshid-ds/android-payload-toolkit/issues)

---

Made with ❤️ for Android devs
