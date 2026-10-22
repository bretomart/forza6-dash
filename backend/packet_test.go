package main

import (
	"encoding/binary"
	"encoding/json"
	"math"
	"testing"
)

// TestMotorsportPacketDecode verifies Motorsport specification decoding
func TestMotorsportPacketDecode(t *testing.T) {
	buf := make([]byte, 324) // Standard 324-byte packet (311 payload + 13 padding)

	binary.LittleEndian.PutUint32(buf[0:4], 1)           // IsRaceOn
	binary.LittleEndian.PutUint32(buf[4:8], 250000)      // TimestampMS
	binary.LittleEndian.PutUint32(buf[8:12], math.Float32bits(8500.0))  // EngineMaxRpm
	binary.LittleEndian.PutUint32(buf[12:16], math.Float32bits(950.0))  // EngineIdleRpm
	binary.LittleEndian.PutUint32(buf[16:20], math.Float32bits(6200.0)) // CurrentEngineRpm

	// In unshifted Motorsport: Position is 232..244, Speed is 244..248, Power is 248..252, Torque is 252..256
	binary.LittleEndian.PutUint32(buf[244:248], math.Float32bits(55.5)) // Speed in m/s (approx 200 km/h)
	binary.LittleEndian.PutUint32(buf[248:252], math.Float32bits(420000.0)) // Power in Watts
	binary.LittleEndian.PutUint32(buf[252:256], math.Float32bits(580.0))    // Torque in N-m

	// In unshifted Motorsport: Fuel is at 276. Set to invalid horizon fuel (e.g. 500) to force Motorsport fallback
	binary.LittleEndian.PutUint32(buf[288:292], math.Float32bits(999.0)) // Invalid Horizon fuel

	buf[303] = 255 // Full accel
	buf[304] = 0   // No brake
	buf[307] = 5   // Gear 5

	packet, err := parsePayload(buf)
	if err != nil {
		t.Fatalf("Failed to parse motorsport payload: %v", err)
	}

	if packet.IsRaceOn != 1 {
		t.Errorf("Expected IsRaceOn 1, got %d", packet.IsRaceOn)
	}
	if packet.CurrentEngineRpm < 6199 || packet.CurrentEngineRpm > 6201 {
		t.Errorf("Expected CurrentEngineRpm ~6200, got %f", packet.CurrentEngineRpm)
	}
	if packet.Speed < 55.4 || packet.Speed > 55.6 {
		t.Errorf("Expected Speed ~55.5 m/s, got %f", packet.Speed)
	}
	if packet.Gear != 5 {
		t.Errorf("Expected Gear 5, got %d", packet.Gear)
	}
}

// TestHorizonPacketWith12BytePadding verifies Horizon 4 / Horizon 5 324-byte packet with gap
func TestHorizonPacketWith12BytePadding(t *testing.T) {
	buf := make([]byte, 324)

	binary.LittleEndian.PutUint32(buf[0:4], 1)
	binary.LittleEndian.PutUint32(buf[4:8], 999999)
	binary.LittleEndian.PutUint32(buf[8:12], math.Float32bits(9000.0))
	binary.LittleEndian.PutUint32(buf[16:20], math.Float32bits(7200.0))

	// In Horizon: Speed is offset 256, Power is 260, Torque is 264, TireTemp is 268..283, Fuel is 288
	binary.LittleEndian.PutUint32(buf[256:260], math.Float32bits(60.0)) // 60 m/s (~216 km/h)
	binary.LittleEndian.PutUint32(buf[260:264], math.Float32bits(500000.0)) // 500 kW
	binary.LittleEndian.PutUint32(buf[264:268], math.Float32bits(700.0)) // 700 N-m

	// Realistic Horizon Tire Temps in °F: FL=195, FR=196, RL=202, RR=203
	binary.LittleEndian.PutUint32(buf[268:272], math.Float32bits(195.0))
	binary.LittleEndian.PutUint32(buf[272:276], math.Float32bits(196.0))
	binary.LittleEndian.PutUint32(buf[276:280], math.Float32bits(202.0))
	binary.LittleEndian.PutUint32(buf[280:284], math.Float32bits(203.0))

	// Horizon Fuel at 288: 0.85 (85%)
	binary.LittleEndian.PutUint32(buf[288:292], math.Float32bits(0.85))

	// Horizon Gear at 319: 6th gear
	buf[315] = 220 // Accel
	buf[316] = 30  // Brake
	buf[319] = 6   // Gear 6

	packet, err := parsePayload(buf)
	if err != nil {
		t.Fatalf("Failed to parse Horizon packet: %v", err)
	}

	if packet.Speed < 59.9 || packet.Speed > 60.1 {
		t.Errorf("Horizon Speed mismatch: expected ~60.0, got %f", packet.Speed)
	}
	if packet.Fuel < 0.84 || packet.Fuel > 0.86 {
		t.Errorf("Horizon Fuel mismatch: expected ~0.85, got %f", packet.Fuel)
	}
	if packet.TireTemp[0] < 194.0 || packet.TireTemp[0] > 196.0 {
		t.Errorf("Horizon TireTemp[0] mismatch: expected ~195.0, got %f", packet.TireTemp[0])
	}
	if packet.Gear != 6 {
		t.Errorf("Horizon Gear mismatch: expected 6, got %d", packet.Gear)
	}

	// Verify JSON serialization round-trip
	data, err := json.Marshal(packet)
	if err != nil {
		t.Fatalf("JSON serialization failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("JSON output unexpectedly empty")
	}
}

// TestCorruptedOrShortPacket verifies resilience to truncated frames
func TestCorruptedOrShortPacket(t *testing.T) {
	shortBuf := make([]byte, 50)
	_, err := parsePayload(shortBuf)
	if err == nil {
		t.Errorf("Expected error when parsing undersized buffer of 50 bytes")
	}
}
