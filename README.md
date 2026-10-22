# Forza 6 Telemetry Dash (Cyber GT Edition)

A high-performance, real-time CAN-bus telemetry HUD and relay server engineered for Forza Horizon and Forza Motorsport titles. It captures incoming 60 Hz UDP Little-Endian packets, decodes vehicle dynamics (speed, RPM, gear, pedals, G-forces, drift slip angle, tire thermals, and suspension travel), and renders a responsive dashboard.

## Standalone Zero-Configuration Releases

Both distributions embed the complete backend relay daemon and the precompiled Single Page Application (SPA) frontend directly into a single self-contained executable. Launching either executable starts the UDP receiver on port 2048, launches the local HTTP/WebSocket bridge, and opens the dashboard directly in the user's default web browser.

| Platform | Binary Artifact | Target Environment | Size |
| :--- | :--- | :--- | :--- |
| **Linux (x86_64)** | `forza6-dash-v1.1.0-linux-x86_64.AppImage` | Universal Linux Distribution | ~4.2 MB |
| **Windows (x64)** | `forza6-dash-v1.1.0-windows-x64.exe` | Windows 11 / Windows 10 (64-bit) | ~7.3 MB |

### Supported Linux Distributions for AppImage
The AppImage package is built using Squashfs with an embedded runtime and functions out of the box across modern Linux distributions:
- **Fedora**
- **Ubuntu**
- **Arch Linux / Manjaro**
- **Debian GNU/Linux**
- **openSUSE**
- **Pop!_OS / Linux Mint / Elementary OS**
- **Red Hat Enterprise Linux / AlmaLinux / Rocky Linux**

*(Note: On systems without FUSE installed, pass `--appimage-extract` to extract and run directly).*

### Running on Windows 11
Double-click `forza6-dash.exe`. Windows Defender SmartScreen may ask for permission on initial launch because the binary is unsigned. Click **More info -> Run anyway**. The application initializes the UDP socket and opens the dashboard in your default browser (Chrome, Edge, Firefox, Brave).

---

## In-Game Configuration (Forza Horizon 4 / 5 / Motorsport)
1. Launch Forza and open **Settings** -> **HUD and Gameplay**.
2. Scroll to the bottom to locate **Data Out**.
3. Configure the telemetry options:
   - **Data Out**: `ON`
   - **Data Out IP Address**: `127.0.0.1` (or your PC's local LAN IP if playing on Xbox/another PC)
   - **Data Out IP Port**: `2048`
   - **Data Out Packet Format**: `Car Dash`

---

## Test Suites & Verification

### Backend Verification (Go Relay)
```bash
cd backend
go test -v ./...
```
- Validates 324-byte packet decoding for standard Motorsport specification.
- Validates Forza Horizon 12-byte padding gap handling (speed at byte 256, fuel at byte 288, tire temps at bytes 268..283).
- Validates error rejection for malformed or truncated UDP packets.

### Frontend Unit & Physics Verification (Bun Test)
```bash
cd frontend
bun test
```
- Validates m/s to metric km/h conversions.
- Validates tire temperatures with 0°C idle guard.
- Validates 6-axis G-force resolution and drift angle vectoring.
- Validates tachometer sweep calculations starting strictly at 0 RPM.
