import { ref, computed } from 'vue'
import { useWebSocket } from '@vueuse/core'
import type { TelemetryPacket } from '../types/telemetry'

export const defaultTelemetryPacket: TelemetryPacket = {
  isRaceOn: 0,
  timestampMS: 0,
  engineMaxRpm: 8000,
  engineIdleRpm: 1000,
  currentEngineRpm: 0,
  accelerationX: 0,
  accelerationY: 0,
  accelerationZ: 0,
  velocityX: 0,
  velocityY: 0,
  velocityZ: 0,
  angularVelocityX: 0,
  angularVelocityY: 0,
  angularVelocityZ: 0,
  yaw: 0,
  pitch: 0,
  roll: 0,
  normalizedSuspensionTravel: [0, 0, 0, 0],
  tireSlipRatio: [0, 0, 0, 0],
  wheelRotationSpeed: [0, 0, 0, 0],
  wheelOnRumble: [0, 0, 0, 0],
  wheelInPuddleDepth: [0, 0, 0, 0],
  surfaceRumble: [0, 0, 0, 0],
  tireSlipAngle: [0, 0, 0, 0],
  tireCombinedSlip: [0, 0, 0, 0],
  suspensionTravelMeters: [0, 0, 0, 0],
  carOrdinal: 0,
  carClass: 0,
  carPerformanceIndex: 100,
  drivetrainType: 0,
  numCylinders: 4,
  positionX: 0,
  positionY: 0,
  positionZ: 0,
  speed: 0,
  power: 0,
  torque: 0,
  tireTemp: [0, 0, 0, 0],
  boost: 0,
  fuel: 1,
  distanceTraveled: 0,
  bestLap: 0,
  lastLap: 0,
  currentLap: 0,
  currentRaceTime: 0,
  lapNumber: 0,
  racePosition: 0,
  accel: 0,
  brake: 0,
  clutch: 0,
  handbrake: 0,
  gear: 11, // Neutral default display
  steer: 0,
  normalizedDrivingLine: 0,
  normalizedAIBrakeDifference: 0,
}

// Module-level reactive singleton store
const telemetry = ref<TelemetryPacket>({ ...defaultTelemetryPacket })
const lastPacketTimestamp = ref<number>(0)
const isLiveStreamActive = ref<boolean>(false)

let staleTimer: ReturnType<typeof setTimeout> | null = null

const { status, open, close } = useWebSocket('ws://localhost:8080/ws', {
  autoReconnect: {
    retries: Infinity,
    delay: 1500,
  },
  onMessage: (_ws, event) => {
    try {
      const parsed = JSON.parse(event.data) as TelemetryPacket
      telemetry.value = parsed
      lastPacketTimestamp.value = Date.now()
      isLiveStreamActive.value = true

      if (staleTimer) clearTimeout(staleTimer)
      // If no packet received for 1.5 seconds, set stream active to false
      staleTimer = setTimeout(() => {
        isLiveStreamActive.value = false
      }, 1500)
    } catch {
      // Discard invalid JSON frames gracefully
    }
  },
})

// Speed computed properties (packet speed is in m/s)
const speedKmh = computed(() => telemetry.value.speed * 3.6)
const speedMph = computed(() => telemetry.value.speed * 2.23694)

// Gear display computed property ("R", "N", "1" through "10")
const gearDisplay = computed(() => {
  const g = telemetry.value.gear
  if (g === 0) return 'R'
  if (g >= 1 && g <= 10) return String(g)
  return 'N'
})

// Normalized RPM percentage against EngineMaxRpm (clamped between 0 and 100%)
const rpmPercentage = computed(() => {
  const maxRpm = telemetry.value.engineMaxRpm || 8000
  const curRpm = telemetry.value.currentEngineRpm || 0
  const pct = (curRpm / maxRpm) * 100
  return Math.min(Math.max(pct, 0), 100)
})

// Pedal inputs mapped to 0-100%
const accelPercentage = computed(() => Math.min(Math.max((telemetry.value.accel / 255) * 100, 0), 100))
const brakePercentage = computed(() => Math.min(Math.max((telemetry.value.brake / 255) * 100, 0), 100))
const clutchPercentage = computed(() => Math.min(Math.max((telemetry.value.clutch / 255) * 100, 0), 100))
const handbrakePercentage = computed(() => Math.min(Math.max((telemetry.value.handbrake / 255) * 100, 0), 100))

// Steering angle normalized from -127..127 to -100%..100% or degrees
const steerPercentage = computed(() => {
  return Math.min(Math.max((telemetry.value.steer / 127) * 100, -100), 100)
})

// G-force calculations (1 G = 9.80665 m/s^2)
const lateralG = computed(() => {
  return (telemetry.value.accelerationX || 0) / 9.80665
})

const longitudinalG = computed(() => {
  return (telemetry.value.accelerationZ || 0) / 9.80665
})

const verticalG = computed(() => {
  return (telemetry.value.accelerationY || 0) / 9.80665
})

const totalG = computed(() => {
  const lat = lateralG.value
  const lon = longitudinalG.value
  return Math.sqrt(lat * lat + lon * lon)
})

// Drift Angle (Slip Angle of the vehicle in degrees)
// Calculated from lateral velocity vs longitudinal velocity: atan2(velocityX, velocityZ)
const driftAngleDeg = computed(() => {
  const vx = telemetry.value.velocityX || 0
  const vz = telemetry.value.velocityZ || 0
  const speed = telemetry.value.speed || 0
  // Require minimum speed threshold of 3 m/s (~11 km/h) to avoid noise while stopped
  if (speed < 3) return 0
  const angleRad = Math.atan2(vx, vz)
  const angleDeg = angleRad * (180 / Math.PI)
  return Math.min(Math.max(angleDeg, -90), 90)
})

// Drift status flag: active when drift angle is above 8 degrees and speed > 15 km/h
const isDrifting = computed(() => {
  return Math.abs(driftAngleDeg.value) > 8 && (telemetry.value.speed || 0) * 3.6 > 15
})

export function useTelemetryStore() {
  return {
    telemetry,
    status,
    isLiveStreamActive,
    lastPacketTimestamp,
    open,
    close,
    speedKmh,
    speedMph,
    gearDisplay,
    rpmPercentage,
    accelPercentage,
    brakePercentage,
    clutchPercentage,
    handbrakePercentage,
    steerPercentage,
    lateralG,
    longitudinalG,
    verticalG,
    totalG,
    driftAngleDeg,
    isDrifting,
  }
}
