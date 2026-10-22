import { describe, it, expect } from 'bun:test'
import { defaultTelemetryPacket } from '../stores/telemetryStore'

describe('Forza Telemetry Core Mathematics & Computations', () => {
  it('default telemetry packet initializes with metric-safe zero values', () => {
    expect(defaultTelemetryPacket.isRaceOn).toBe(0)
    expect(defaultTelemetryPacket.speed).toBe(0)
    expect(defaultTelemetryPacket.tireTemp).toEqual([0, 0, 0, 0])
    expect(defaultTelemetryPacket.engineMaxRpm).toBe(8000)
    expect(defaultTelemetryPacket.gear).toBe(11) // Default Neutral
  })

  it('correctly converts m/s to metric km/h', () => {
    // 50 m/s = 180 km/h
    const speedMs = 50.0
    const speedKmh = speedMs * 3.6
    expect(Math.round(speedKmh)).toBe(180)

    // 0 m/s = 0 km/h
    expect(0 * 3.6).toBe(0)

    // 83.33 m/s ~= 300 km/h
    expect(Math.round(83.3333 * 3.6)).toBe(300)
  })

  it('correctly maps gear values to automotive cluster notation', () => {
    const formatGear = (g: number) => {
      if (g === 0) return 'R'
      if (g >= 1 && g <= 10) return String(g)
      return 'N'
    }

    expect(formatGear(0)).toBe('R')
    expect(formatGear(1)).toBe('1')
    expect(formatGear(6)).toBe('6')
    expect(formatGear(10)).toBe('10')
    expect(formatGear(11)).toBe('N')
    expect(formatGear(-1)).toBe('N')
  })

  it('computes metric tire temperature in Celsius with 0°C idle guard', () => {
    const calcTempC = (rawF: number) => {
      if (!rawF || rawF <= 0) return 0
      return Math.round((rawF - 32) * (5 / 9))
    }

    // Idle state (packet sends 0) must display 0°C (not -18°C)
    expect(calcTempC(0)).toBe(0)
    expect(calcTempC(undefined as unknown as number)).toBe(0)

    // Freezing: 32°F = 0°C
    expect(calcTempC(32)).toBe(0)

    // Operating temp: 200°F = 93°C
    expect(calcTempC(200)).toBe(93)

    // Hot race tire: 230°F = 110°C
    expect(calcTempC(230)).toBe(110)
  })

  it('computes boost pressure from PSI to Bar correctly', () => {
    const psiToBar = (psi: number) => Math.max(psi * 0.0689476, 0)

    expect(psiToBar(0)).toBe(0)
    // 14.5 psi ~= 1.00 bar
    expect(parseFloat(psiToBar(14.5038).toFixed(2))).toBe(1.00)
    // 29.0 psi ~= 2.00 bar
    expect(parseFloat(psiToBar(29.0075).toFixed(2))).toBe(2.00)
  })

  it('computes powertrain kW and Nm conversions accurately', () => {
    // Forza packet power is in Watts
    const powerWatts = 450000 // 450 kW (~603 hp)
    const powerKw = Math.max(Math.round(powerWatts / 1000), 0)
    expect(powerKw).toBe(450)

    // Forza packet torque is in N-m
    const rawTorque = 650.4
    const torqueNm = Math.max(Math.round(rawTorque), 0)
    expect(torqueNm).toBe(650)
  })

  it('computes lateral and longitudinal G-forces from m/s^2', () => {
    const gConst = 9.80665
    const accelX = 14.709975 // 1.5 G lateral
    const accelZ = -9.80665  // -1.0 G braking deceleration

    const latG = accelX / gConst
    const lonG = accelZ / gConst
    const totalG = Math.sqrt(latG * latG + lonG * lonG)

    expect(parseFloat(latG.toFixed(2))).toBe(1.50)
    expect(parseFloat(lonG.toFixed(2))).toBe(-1.00)
    expect(parseFloat(totalG.toFixed(2))).toBe(1.80)
  })

  it('calculates vehicle drift angle from velocity vectors', () => {
    const calcDrift = (vx: number, vz: number, speed: number) => {
      if (speed < 3) return 0
      const angleRad = Math.atan2(vx, vz)
      const angleDeg = angleRad * (180 / Math.PI)
      return Math.min(Math.max(angleDeg, -90), 90)
    }

    // Straight line: vx = 0, vz = 40 m/s
    expect(calcDrift(0, 40, 40)).toBe(0)

    // Stationary noise rejection
    expect(calcDrift(2, 2, 2.5)).toBe(0)

    // 45 degree slide
    expect(Math.round(calcDrift(30, 30, 42.4))).toBe(45)

    // Negative 30 degree slide
    expect(Math.round(calcDrift(-20, 34.64, 40))).toBe(-30)
  })

  it('calculates tachometer arc sweep strictly starting from 0', () => {
    const arcLength = 485.90
    const calcDashArray = (currentRpm: number, maxRpm: number) => {
      const norm = Math.min(Math.max(currentRpm / maxRpm, 0), 1)
      const activeLength = norm * arcLength
      return { norm, activeLength }
    }

    // At 0 RPM: stroke length must be strictly 0
    const idleZero = calcDashArray(0, 8000)
    expect(idleZero.norm).toBe(0)
    expect(idleZero.activeLength).toBe(0)

    // At 4000 RPM (50%): stroke is half
    const mid = calcDashArray(4000, 8000)
    expect(mid.norm).toBe(0.5)
    expect(mid.activeLength).toBeCloseTo(242.95, 1)

    // At redline (8000 RPM): full 240 deg arc
    const redline = calcDashArray(8000, 8000)
    expect(redline.norm).toBe(1.0)
    expect(redline.activeLength).toBeCloseTo(485.90, 1)
  })

  it('dynamically accumulates session distance from speed and time interval', () => {
    let accumulatedDistanceMeters = 0
    let lastTime = 1000

    const processPacket = (speedMs: number, currentTimeMs: number, packetDistance: number) => {
      if (lastTime > 0 && speedMs > 0.1) {
        const dt = (currentTimeMs - lastTime) / 1000
        if (dt > 0 && dt < 0.5) {
          accumulatedDistanceMeters += speedMs * dt
        }
      }
      lastTime = currentTimeMs

      if (packetDistance && packetDistance > accumulatedDistanceMeters) {
        accumulatedDistanceMeters = packetDistance
      }
      return accumulatedDistanceMeters
    }

    // Driving at 50 m/s (~180 km/h) for 100ms
    const d1 = processPacket(50.0, 1100, 0)
    expect(d1).toBeCloseTo(5.0, 1)

    // Driving another 100ms at 50 m/s
    const d2 = processPacket(50.0, 1200, 0)
    expect(d2).toBeCloseTo(10.0, 1)

    // Official track measurement in an active race supersedes if higher
    const d3 = processPacket(50.0, 1300, 250.0)
    expect(d3).toBe(250.0)

    // Reconnection lag guard: if dt >= 0.5s, don't create artificial jump
    const d4 = processPacket(50.0, 3000, 0)
    expect(d4).toBe(250.0)
  })
})
