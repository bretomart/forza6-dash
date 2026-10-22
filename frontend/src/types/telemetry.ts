export interface TelemetryPacket {
  isRaceOn: number
  timestampMS: number
  engineMaxRpm: number
  engineIdleRpm: number
  currentEngineRpm: number
  accelerationX: number
  accelerationY: number
  accelerationZ: number
  velocityX: number
  velocityY: number
  velocityZ: number
  angularVelocityX: number
  angularVelocityY: number
  angularVelocityZ: number
  yaw: number
  pitch: number
  roll: number
  normalizedSuspensionTravel: [number, number, number, number] // FL, FR, RL, RR
  tireSlipRatio: [number, number, number, number] // FL, FR, RL, RR
  wheelRotationSpeed: [number, number, number, number] // FL, FR, RL, RR
  wheelOnRumble: [number, number, number, number]
  wheelInPuddleDepth: [number, number, number, number]
  surfaceRumble: [number, number, number, number]
  tireSlipAngle: [number, number, number, number]
  tireCombinedSlip: [number, number, number, number]
  suspensionTravelMeters: [number, number, number, number]
  carOrdinal: number
  carClass: number
  carPerformanceIndex: number
  drivetrainType: number
  numCylinders: number
  positionX: number
  positionY: number
  positionZ: number
  speed: number // m/s
  power: number // Watts
  torque: number // N-m
  tireTemp: [number, number, number, number] // FL, FR, RL, RR in Fahrenheit
  boost: number // psi
  fuel: number // 0.0 to 1.0
  distanceTraveled: number
  bestLap: number
  lastLap: number
  currentLap: number
  currentRaceTime: number
  lapNumber: number
  racePosition: number
  accel: number // 0-255
  brake: number // 0-255
  clutch: number // 0-255
  handbrake: number // 0-255
  gear: number // 0 = R, 1-10 = Forward gears
  steer: number // -127 to 127
  normalizedDrivingLine: number
  normalizedAIBrakeDifference: number
}
