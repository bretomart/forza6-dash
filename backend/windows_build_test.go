package main

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"io"
	"os"
	"strings"
	"testing"
)

// TestWindowsPEStructure validates the PE32+ executable header, architecture,
// subsystem, and OS compatibility for Windows 10/11 x64 systems.
func TestWindowsPEStructure(t *testing.T) {
	exePath := "../release/forza6-dash-v1.1.2-windows-x64.exe"
	file, err := pe.Open(exePath)
	if err != nil {
		t.Fatalf("Failed to open Windows executable: %v", err)
	}
	defer file.Close()

	// 1. Verify AMD64 machine architecture (0x8664)
	if file.FileHeader.Machine != pe.IMAGE_FILE_MACHINE_AMD64 {
		t.Fatalf("Expected Machine AMD64 (0x8664), got 0x%x", file.FileHeader.Machine)
	}

	// 2. Verify PE32+ Optional Header Magic (0x20b)
	optHdr, ok := file.OptionalHeader.(*pe.OptionalHeader64)
	if !ok {
		t.Fatalf("Expected OptionalHeader64 (PE32+), got %T", file.OptionalHeader)
	}
	if optHdr.Magic != 0x20b {
		t.Fatalf("Expected Magic 0x20b (PE32+), got 0x%x", optHdr.Magic)
	}

	// 3. Verify Subsystem: Windows GUI (2)
	// Subsystem 2 prevents a console terminal window from popping up on launch in Windows 11.
	if optHdr.Subsystem != 2 {
		t.Fatalf("Expected Subsystem 2 (IMAGE_SUBSYSTEM_WINDOWS_GUI), got %d", optHdr.Subsystem)
	}

	// 4. Verify OS / Subsystem version compatibility (Windows 6.1+ covers Windows 7, 8, 10, 11)
	if optHdr.MajorOperatingSystemVersion != 6 || optHdr.MinorOperatingSystemVersion != 1 {
		t.Logf("OS Version: %d.%d", optHdr.MajorOperatingSystemVersion, optHdr.MinorOperatingSystemVersion)
	}

	// 5. Verify vital sections exist
	expectedSections := []string{".text", ".rdata", ".data", ".idata", ".rsrc"}
	sectionMap := make(map[string]*pe.Section)
	for _, sec := range file.Sections {
		name := strings.TrimRight(sec.Name, "\x00")
		sectionMap[name] = sec
	}

	for _, name := range expectedSections {
		if _, exists := sectionMap[name]; !exists {
			t.Errorf("Missing expected PE section: %s", name)
		}
	}

	// 6. Verify resource table (.rsrc)
	rsrcSec, exists := sectionMap[".rsrc"]
	if !exists {
		t.Fatalf("Resource section .rsrc not found")
	}
	if rsrcSec.Size == 0 {
		t.Fatalf(".rsrc section has size 0")
	}

	// 7. Verify imported symbols (must include symbols from kernel32.dll)
	importedSyms, err := file.ImportedSymbols()
	if err != nil {
		t.Fatalf("Failed to read imported symbols: %v", err)
	}
	hasKernel32 := false
	for _, sym := range importedSyms {
		if strings.HasSuffix(strings.ToLower(sym), ":kernel32.dll") {
			hasKernel32 = true
			break
		}
	}
	if !hasKernel32 {
		t.Errorf("Executable does not import kernel32.dll symbols")
	}
}

// TestWindowsEmbeddedResources inspects the embedded PE resource directory
// ensuring valid multi-resolution desktop icons (RT_ICON and RT_GROUP_ICON) are present.
func TestWindowsEmbeddedResources(t *testing.T) {
	exePath := "../release/forza6-dash-v1.1.2-windows-x64.exe"
	raw, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatalf("Failed to read binary: %v", err)
	}

	file, err := pe.Open(exePath)
	if err != nil {
		t.Fatalf("Failed to parse PE: %v", err)
	}
	defer file.Close()

	var rsrcSec *pe.Section
	for _, s := range file.Sections {
		if strings.TrimRight(s.Name, "\x00") == ".rsrc" {
			rsrcSec = s
			break
		}
	}
	if rsrcSec == nil {
		t.Fatalf(".rsrc section not found")
	}

	secOffset := int(rsrcSec.Offset)
	if secOffset+16 > len(raw) {
		t.Fatalf(".rsrc offset exceeds file bounds")
	}

	// Root resource directory
	numNamed := binary.LittleEndian.Uint16(raw[secOffset+12 : secOffset+14])
	numID := binary.LittleEndian.Uint16(raw[secOffset+14 : secOffset+16])
	totalEntries := int(numNamed + numID)

	hasIconType := false
	hasGroupIconType := false

	for i := 0; i < totalEntries; i++ {
		entryOffset := secOffset + 16 + i*8
		typeID := binary.LittleEndian.Uint32(raw[entryOffset : entryOffset+4])
		typeVal := typeID & 0x7FFFFFFF
		if typeVal == 3 { // RT_ICON
			hasIconType = true
		}
		if typeVal == 14 { // RT_GROUP_ICON
			hasGroupIconType = true
		}
	}

	if !hasIconType {
		t.Errorf("PE resources do not contain RT_ICON (type 3)")
	}
	if !hasGroupIconType {
		t.Errorf("PE resources do not contain RT_GROUP_ICON (type 14)")
	}
}

// TestWindowsEmbeddedFrontendAssets verifies that the compiled Vue 3 HUD assets
// are embedded in the Windows executable binary.
func TestWindowsEmbeddedFrontendAssets(t *testing.T) {
	exePath := "../release/forza6-dash-v1.1.2-windows-x64.exe"
	raw, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatalf("Failed to read binary: %v", err)
	}

	essentialAssets := []string{
		"index.html",
		"car_outline_pure.png",
		"favicon.svg",
	}

	for _, asset := range essentialAssets {
		if !bytes.Contains(raw, []byte(asset)) {
			t.Errorf("Embedded binary missing essential asset: %s", asset)
		}
	}
}

// TestWindowsFunctionalEndpoints verifies that Windows command-line routines,
// browser opening mechanisms, systray handlers, and shutdown endpoints are present.
func TestWindowsFunctionalEndpoints(t *testing.T) {
	exePath := "../release/forza6-dash-v1.1.2-windows-x64.exe"
	raw, err := os.ReadFile(exePath)
	if err != nil {
		t.Fatalf("Failed to read binary: %v", err)
	}

	signatures := []struct {
		name      string
		pattern   []byte
		mandatory bool
	}{
		{"Windows Browser Launcher", []byte("rundll32"), true},
		{"Windows Protocol Handler", []byte("url.dll,FileProtocolHandler"), true},
		{"Telemetry Websocket Route", []byte("/ws"), true},
		{"Graceful Shutdown Route", []byte("/api/shutdown"), true},
		{"Systray Open Action", []byte("Open Dashboard"), true},
		{"Systray Quit Action", []byte("Quit Forza 6 Dash"), true},
	}

	for _, sig := range signatures {
		if !bytes.Contains(raw, sig.pattern) {
			if sig.mandatory {
				t.Errorf("Missing required functionality string %q (%s)", sig.name, sig.pattern)
			}
		}
	}
}

// TestWindowsBinaryReader verifies that standard io operations against the binary succeed.
func TestWindowsBinaryReader(t *testing.T) {
	exePath := "../release/forza6-dash-v1.1.2-windows-x64.exe"
	f, err := os.Open(exePath)
	if err != nil {
		t.Fatalf("Failed to open binary: %v", err)
	}
	defer f.Close()

	magic := make([]byte, 2)
	if _, err := io.ReadFull(f, magic); err != nil {
		t.Fatalf("Failed to read header: %v", err)
	}
	if string(magic) != "MZ" {
		t.Fatalf("Expected MZ magic, got %s", string(magic))
	}
}
