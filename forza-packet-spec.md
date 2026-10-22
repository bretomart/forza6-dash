# Forza Data Out (Dash Format) Specification

Total Packet Size: 324 bytes (Little Endian)

## Core Offsets and Types
- 0: s32 IsRaceOn
- 4: u32 TimestampMS
- 8: f32 EngineMaxRpm
- 12: f32 EngineIdleRpm
- 16: f32 CurrentEngineRpm
- 20-31: f32 AccelerationX, AccelerationY, AccelerationZ
- 32-43: f32 VelocityX, VelocityY, VelocityZ
- 44-55: f32 AngularVelocityX, AngularVelocityY, AngularVelocityZ
- 56-67: f32 Yaw, Pitch, Roll
- 68-83: f32 NormalizedSuspensionTravel (FL, FR, RL, RR)
- 84-99: f32 TireSlipRatio (FL, FR, RL, RR)
- 100-115: f32 WheelRotationSpeed (FL, FR, RL, RR)
- 244: f32 Speed (m/s, multiply by 3.6 for km/h or 2.23694 for mph)
- 248: f32 Power (Watts)
- 252: f32 Torque (Newton-meters)
- 256-271: f32 TireTemp (FL, FR, RL, RR in Fahrenheit)
- 272: f32 Boost (psi)
- 276: f32 Fuel (0.0 to 1.0)
- 300: u16 LapNumber
- 302: u8 RacePosition
- 303: u8 Accel (0-255)
- 304: u8 Brake (0-255)
- 305: u8 Clutch (0-255)
- 306: u8 Handbrake (0-255)
- 307: u8 Gear (0 = Reverse, 1-10 = Forward gears)
- 308: s8 Steer (-127 to 127)
