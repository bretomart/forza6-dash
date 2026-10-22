package main

import (
	"bytes"
	"embed"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//go:embed web/*
var embeddedDashboard embed.FS

// TelemetryPacket maps the Forza telemetry payload for both Motorsport and Horizon layouts.
type TelemetryPacket struct {
	IsRaceOn           int32      `json:"isRaceOn"`
	TimestampMS        uint32     `json:"timestampMS"`
	EngineMaxRpm       float32    `json:"engineMaxRpm"`
	EngineIdleRpm      float32    `json:"engineIdleRpm"`
	CurrentEngineRpm   float32    `json:"currentEngineRpm"`
	AccelerationX      float32    `json:"accelerationX"`
	AccelerationY      float32    `json:"accelerationY"`
	AccelerationZ      float32    `json:"accelerationZ"`
	VelocityX          float32    `json:"velocityX"`
	VelocityY          float32    `json:"velocityY"`
	VelocityZ          float32    `json:"velocityZ"`
	AngularVelocityX   float32    `json:"angularVelocityX"`
	AngularVelocityY   float32    `json:"angularVelocityY"`
	AngularVelocityZ   float32    `json:"angularVelocityZ"`
	Yaw                float32    `json:"yaw"`
	Pitch              float32    `json:"pitch"`
	Roll               float32    `json:"roll"`
	NormalizedSuspTravel [4]float32 `json:"normalizedSuspensionTravel"`
	TireSlipRatio      [4]float32 `json:"tireSlipRatio"`
	WheelRotationSpeed [4]float32 `json:"wheelRotationSpeed"`
	WheelOnRumble      [4]int32   `json:"wheelOnRumble"`
	WheelInPuddleDepth [4]float32 `json:"wheelInPuddleDepth"`
	SurfaceRumble      [4]float32 `json:"surfaceRumble"`
	TireSlipAngle      [4]float32 `json:"tireSlipAngle"`
	TireCombinedSlip   [4]float32 `json:"tireCombinedSlip"`
	SuspensionTravelMeters [4]float32 `json:"suspensionTravelMeters"`
	CarOrdinal         int32      `json:"carOrdinal"`
	CarClass           int32      `json:"carClass"`
	CarPerformanceIndex int32     `json:"carPerformanceIndex"`
	DrivetrainType     int32      `json:"drivetrainType"`
	NumCylinders       int32      `json:"numCylinders"`
	PositionX          float32    `json:"positionX"`
	PositionY          float32    `json:"positionY"`
	PositionZ          float32    `json:"positionZ"`
	Speed              float32    `json:"speed"`
	Power              float32    `json:"power"`
	Torque             float32    `json:"torque"`
	TireTemp           [4]float32 `json:"tireTemp"`
	Boost              float32    `json:"boost"`
	Fuel               float32    `json:"fuel"`
	DistanceTraveled   float32    `json:"distanceTraveled"`
	BestLap            float32    `json:"bestLap"`
	LastLap            float32    `json:"lastLap"`
	CurrentLap         float32    `json:"currentLap"`
	CurrentRaceTime    float32    `json:"currentRaceTime"`
	LapNumber          uint16     `json:"lapNumber"`
	RacePosition       uint8      `json:"racePosition"`
	Accel              uint8      `json:"accel"`
	Brake              uint8      `json:"brake"`
	Clutch             uint8      `json:"clutch"`
	Handbrake          uint8      `json:"handbrake"`
	Gear               uint8      `json:"gear"`
	Steer              int8       `json:"steer"`
	NormalizedDrivingLine int8    `json:"normalizedDrivingLine"`
	NormalizedAIBrakeDifference int8 `json:"normalizedAIBrakeDifference"`
}

// RawHorizonPacket represents the 324-byte packet layout with the Horizon 12-byte padding gap at bytes 232..243.
type RawHorizonPacket struct {
	IsRaceOn           int32
	TimestampMS        uint32
	EngineMaxRpm       float32
	EngineIdleRpm      float32
	CurrentEngineRpm   float32
	AccelerationX      float32
	AccelerationY      float32
	AccelerationZ      float32
	VelocityX          float32
	VelocityY          float32
	VelocityZ          float32
	AngularVelocityX   float32
	AngularVelocityY   float32
	AngularVelocityZ   float32
	Yaw                float32
	Pitch              float32
	Roll               float32
	NormalizedSuspTravel [4]float32
	TireSlipRatio      [4]float32
	WheelRotationSpeed [4]float32
	WheelOnRumble      [4]int32
	WheelInPuddleDepth [4]float32
	SurfaceRumble      [4]float32
	TireSlipAngle      [4]float32
	TireCombinedSlip   [4]float32
	SuspensionTravelMeters [4]float32
	CarOrdinal         int32
	CarClass           int32
	CarPerformanceIndex int32
	DrivetrainType     int32
	NumCylinders       int32
	HorizonPadding     [12]byte   // bytes 232..243
	PositionX          float32    // 244
	PositionY          float32    // 248
	PositionZ          float32    // 252
	Speed              float32    // 256
	Power              float32    // 260
	Torque             float32    // 264
	TireTemp           [4]float32 // 268..283
	Boost              float32    // 284
	Fuel               float32    // 288
	DistanceTraveled   float32    // 292
	BestLap            float32    // 296
	LastLap            float32    // 300
	CurrentLap         float32    // 304
	CurrentRaceTime    float32    // 308
	LapNumber          uint16     // 312
	RacePosition       uint8      // 314
	Accel              uint8      // 315
	Brake              uint8      // 316
	Clutch             uint8      // 317
	Handbrake          uint8      // 318
	Gear               uint8      // 319
	Steer              int8       // 320
	NormalizedDrivingLine int8    // 321
	NormalizedAIBrakeDiff int8    // 322
	ReservedPadding    [1]byte    // 323 (total 324 bytes)
}

// RawMotorsportPacket represents the FM7 311-byte or unshifted 324-byte spec layout.
type RawMotorsportPacket struct {
	IsRaceOn           int32
	TimestampMS        uint32
	EngineMaxRpm       float32
	EngineIdleRpm      float32
	CurrentEngineRpm   float32
	AccelerationX      float32
	AccelerationY      float32
	AccelerationZ      float32
	VelocityX          float32
	VelocityY          float32
	VelocityZ          float32
	AngularVelocityX   float32
	AngularVelocityY   float32
	AngularVelocityZ   float32
	Yaw                float32
	Pitch              float32
	Roll               float32
	NormalizedSuspTravel [4]float32
	TireSlipRatio      [4]float32
	WheelRotationSpeed [4]float32
	WheelOnRumble      [4]int32
	WheelInPuddleDepth [4]float32
	SurfaceRumble      [4]float32
	TireSlipAngle      [4]float32
	TireCombinedSlip   [4]float32
	SuspensionTravelMeters [4]float32
	CarOrdinal         int32
	CarClass           int32
	CarPerformanceIndex int32
	DrivetrainType     int32
	NumCylinders       int32
	PositionX          float32
	PositionY          float32
	PositionZ          float32
	Speed              float32
	Power              float32
	Torque             float32
	TireTemp           [4]float32
	Boost              float32
	Fuel               float32
	DistanceTraveled   float32
	BestLap            float32
	LastLap            float32
	CurrentLap         float32
	CurrentRaceTime    float32
	LapNumber          uint16
	RacePosition       uint8
	Accel              uint8
	Brake              uint8
	Clutch             uint8
	Handbrake          uint8
	Gear               uint8
	Steer              int8
	NormalizedDrivingLine int8
	NormalizedAIBrakeDiff int8
	ReservedPadding    [13]byte
}

func parsePayload(buf []byte) (TelemetryPacket, error) {
	// If 324 bytes, test whether it matches Horizon layout (speed at 256, fuel at 288)
	// Horizon fuel at 288 is typically 0.0 .. 1.0, and tire temps at 268..283 are between 50 and 400 °F
	if len(buf) == 324 {
		var h RawHorizonPacket
		reader := bytes.NewReader(buf)
		if err := binary.Read(reader, binary.LittleEndian, &h); err == nil {
			// Check if Horizon offsets produce realistic tire temps (50°F..400°F) or realistic fuel (0..1.05)
			if (h.TireTemp[0] >= 30 && h.TireTemp[0] <= 500) || (h.Fuel >= 0 && h.Fuel <= 1.05) {
				return TelemetryPacket{
					IsRaceOn:           h.IsRaceOn,
					TimestampMS:        h.TimestampMS,
					EngineMaxRpm:       h.EngineMaxRpm,
					EngineIdleRpm:      h.EngineIdleRpm,
					CurrentEngineRpm:   h.CurrentEngineRpm,
					AccelerationX:      h.AccelerationX,
					AccelerationY:      h.AccelerationY,
					AccelerationZ:      h.AccelerationZ,
					VelocityX:          h.VelocityX,
					VelocityY:          h.VelocityY,
					VelocityZ:          h.VelocityZ,
					AngularVelocityX:   h.AngularVelocityX,
					AngularVelocityY:   h.AngularVelocityY,
					AngularVelocityZ:   h.AngularVelocityZ,
					Yaw:                h.Yaw,
					Pitch:              h.Pitch,
					Roll:               h.Roll,
					NormalizedSuspTravel: h.NormalizedSuspTravel,
					TireSlipRatio:      h.TireSlipRatio,
					WheelRotationSpeed: h.WheelRotationSpeed,
					WheelOnRumble:      h.WheelOnRumble,
					WheelInPuddleDepth: h.WheelInPuddleDepth,
					SurfaceRumble:      h.SurfaceRumble,
					TireSlipAngle:      h.TireSlipAngle,
					TireCombinedSlip:   h.TireCombinedSlip,
					SuspensionTravelMeters: h.SuspensionTravelMeters,
					CarOrdinal:         h.CarOrdinal,
					CarClass:           h.CarClass,
					CarPerformanceIndex: h.CarPerformanceIndex,
					DrivetrainType:     h.DrivetrainType,
					NumCylinders:       h.NumCylinders,
					PositionX:          h.PositionX,
					PositionY:          h.PositionY,
					PositionZ:          h.PositionZ,
					Speed:              h.Speed,
					Power:              h.Power,
					Torque:             h.Torque,
					TireTemp:           h.TireTemp,
					Boost:              h.Boost,
					Fuel:               h.Fuel,
					DistanceTraveled:   h.DistanceTraveled,
					BestLap:            h.BestLap,
					LastLap:            h.LastLap,
					CurrentLap:         h.CurrentLap,
					CurrentRaceTime:    h.CurrentRaceTime,
					LapNumber:          h.LapNumber,
					RacePosition:       h.RacePosition,
					Accel:              h.Accel,
					Brake:              h.Brake,
					Clutch:             h.Clutch,
					Handbrake:          h.Handbrake,
					Gear:               h.Gear,
					Steer:              h.Steer,
					NormalizedDrivingLine: h.NormalizedDrivingLine,
					NormalizedAIBrakeDifference: h.NormalizedAIBrakeDiff,
				}, nil
			}
		}
	}

	// Fallback to Motorsport spec layout
	var m RawMotorsportPacket
	reader := bytes.NewReader(buf)
	if err := binary.Read(reader, binary.LittleEndian, &m); err != nil {
		return TelemetryPacket{}, err
	}
	return TelemetryPacket{
		IsRaceOn:           m.IsRaceOn,
		TimestampMS:        m.TimestampMS,
		EngineMaxRpm:       m.EngineMaxRpm,
		EngineIdleRpm:      m.EngineIdleRpm,
		CurrentEngineRpm:   m.CurrentEngineRpm,
		AccelerationX:      m.AccelerationX,
		AccelerationY:      m.AccelerationY,
		AccelerationZ:      m.AccelerationZ,
		VelocityX:          m.VelocityX,
		VelocityY:          m.VelocityY,
		VelocityZ:          m.VelocityZ,
		AngularVelocityX:   m.AngularVelocityX,
		AngularVelocityY:   m.AngularVelocityY,
		AngularVelocityZ:   m.AngularVelocityZ,
		Yaw:                m.Yaw,
		Pitch:              m.Pitch,
		Roll:               m.Roll,
		NormalizedSuspTravel: m.NormalizedSuspTravel,
		TireSlipRatio:      m.TireSlipRatio,
		WheelRotationSpeed: m.WheelRotationSpeed,
		WheelOnRumble:      m.WheelOnRumble,
		WheelInPuddleDepth: m.WheelInPuddleDepth,
		SurfaceRumble:      m.SurfaceRumble,
		TireSlipAngle:      m.TireSlipAngle,
		TireCombinedSlip:   m.TireCombinedSlip,
		SuspensionTravelMeters: m.SuspensionTravelMeters,
		CarOrdinal:         m.CarOrdinal,
		CarClass:           m.CarClass,
		CarPerformanceIndex: m.CarPerformanceIndex,
		DrivetrainType:     m.DrivetrainType,
		NumCylinders:       m.NumCylinders,
		PositionX:          m.PositionX,
		PositionY:          m.PositionY,
		PositionZ:          m.PositionZ,
		Speed:              m.Speed,
		Power:              m.Power,
		Torque:             m.Torque,
		TireTemp:           m.TireTemp,
		Boost:              m.Boost,
		Fuel:               m.Fuel,
		DistanceTraveled:   m.DistanceTraveled,
		BestLap:            m.BestLap,
		LastLap:            m.LastLap,
		CurrentLap:         m.CurrentLap,
		CurrentRaceTime:    m.CurrentRaceTime,
		LapNumber:          m.LapNumber,
		RacePosition:       m.RacePosition,
		Accel:              m.Accel,
		Brake:              m.Brake,
		Clutch:             m.Clutch,
		Handbrake:          m.Handbrake,
		Gear:               m.Gear,
		Steer:              m.Steer,
		NormalizedDrivingLine: m.NormalizedDrivingLine,
		NormalizedAIBrakeDifference: m.NormalizedAIBrakeDiff,
	}, nil
}

// Hub manages active WebSocket client connections and packet broadcasts.
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.RWMutex
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte, 128),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			h.mu.Lock()
			h.clients[conn] = true
			h.mu.Unlock()
			log.Printf("[WS] Client connected (%d active)", len(h.clients))
		case conn := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[conn]; ok {
				delete(h.clients, conn)
				conn.Close()
			}
			h.mu.Unlock()
			log.Printf("[WS] Client disconnected (%d active)", len(h.clients))
		case message := <-h.broadcast:
			h.mu.RLock()
			for conn := range h.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("[WS] Write error: %v", err)
					conn.Close()
					delete(h.clients, conn)
				}
			}
			h.mu.RUnlock()
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // Linux, BSD
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func main() {
	udpAddr := flag.String("udp", ":2048", "UDP address to listen for Forza telemetry")
	httpAddr := flag.String("http", ":8080", "HTTP address to serve WebSocket and Dashboard")
	noBrowser := flag.Bool("no-browser", false, "Disable automatically opening web browser on startup")
	flag.Parse()

	hub := newHub()
	go hub.run()

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("[WS] Upgrade error: %v", err)
			return
		}
		hub.register <- conn

		go func() {
			defer func() {
				hub.unregister <- conn
			}()
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()
	})

	// Serve the embedded Vue 3 SPA frontend directly from root
	webSubFS, err := fs.Sub(embeddedDashboard, "web")
	if err == nil {
		fileServer := http.FileServer(http.FS(webSubFS))
		http.Handle("/", fileServer)
	}

	go func() {
		log.Printf("[HTTP] Unified Forza Telemetry Dashboard & Relay online at http://localhost%s", *httpAddr)
		if err := http.ListenAndServe(*httpAddr, nil); err != nil {
			log.Fatalf("[HTTP] Fatal server error: %v", err)
		}
	}()

	// Open user's default browser after a brief delay so server has bound
	if !*noBrowser {
		go func() {
			time.Sleep(400 * time.Millisecond)
			targetURL := fmt.Sprintf("http://localhost%s", *httpAddr)
			if err := openBrowser(targetURL); err != nil {
				log.Printf("[BROWSER] Notice: could not automatically open default browser: %v", err)
			} else {
				log.Printf("[BROWSER] Opened dashboard at %s in default system browser", targetURL)
			}
		}()
	}

	udpConn, err := net.ListenPacket("udp4", *udpAddr)
	if err != nil {
		log.Fatalf("[UDP] Failed to listen on %s: %v", *udpAddr, err)
	}
	defer udpConn.Close()
	log.Printf("[UDP] Listening for Forza packets on %s (UDP4 wildcard: 127.0.0.1 and LAN)", *udpAddr)

	buf := make([]byte, 1024)
	const expectedSize = 324
	var packetCount uint64

	for {
		n, addr, err := udpConn.ReadFrom(buf)
		if err != nil {
			log.Printf("[UDP] Read error: %v", err)
			continue
		}

		if n != expectedSize && n != 311 {
			log.Printf("[UDP] Discarding packet from %v: expected %d bytes, got %d bytes", addr, expectedSize, n)
			continue
		}

		packet, err := parsePayload(buf[:n])
		if err != nil {
			log.Printf("[UDP] Decode error: %v", err)
			continue
		}

		packetCount++
		if packetCount == 1 || (packetCount%300 == 0 && packet.IsRaceOn != 0) {
			log.Printf("[UDP] Live telemetry from %v (total: %d) | race=%d rpm=%.0f speed=%.1f km/h gear=%d FL_temp=%.0f°F fuel=%.0f%%",
				addr, packetCount, packet.IsRaceOn, packet.CurrentEngineRpm, packet.Speed*3.6, packet.Gear, packet.TireTemp[0], packet.Fuel*100)
		}

		jsonData, err := json.Marshal(packet)
		if err != nil {
			log.Printf("[JSON] Marshal error: %v", err)
			continue
		}

		select {
		case hub.broadcast <- jsonData:
		default:
		}
	}
}
