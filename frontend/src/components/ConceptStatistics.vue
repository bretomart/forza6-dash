<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { telemetry, lateralG, longitudinalG, totalG, driftAngleDeg, isDrifting } = useTelemetryStore()

// Real-time rolling history buffer of live multi-metric data points (Speed, RPM, G-Force)
const bufferLen = 35
const speedHistory = ref<number[]>(Array(bufferLen).fill(0))
const rpmHistory = ref<number[]>(Array(bufferLen).fill(0))
const gHistory = ref<number[]>(Array(bufferLen).fill(0))

// Record live metrics every 100ms
const interval = setInterval(() => {
  const spd = (telemetry.value.speed || 0) * 3.6 // 0..300 km/h
  const rpm = (telemetry.value.currentEngineRpm || 0) / (telemetry.value.engineMaxRpm || 8000) // 0..1
  const g = totalG.value // 0..2.0 G

  speedHistory.value.push(Math.min(Math.max((spd / 250) * 50, 4), 55))
  rpmHistory.value.push(Math.min(Math.max(rpm * 50, 4), 55))
  gHistory.value.push(Math.min(Math.max((g / 1.5) * 50, 4), 55))

  if (speedHistory.value.length > bufferLen) speedHistory.value.shift()
  if (rpmHistory.value.length > bufferLen) rpmHistory.value.shift()
  if (gHistory.value.length > bufferLen) gHistory.value.shift()
}, 100)

onUnmounted(() => {
  clearInterval(interval)
})

// Build smooth SVG Catmull-Rom spline curves from the live buffers
const createSplinePath = (data: number[], height: number) => {
  const width = 300
  const step = width / (data.length - 1)
  const points = data.map((val, i) => ({
    x: i * step,
    y: height - val
  }))

  if (points.length === 0) return ''

  let d = `M ${points[0].x} ${points[0].y}`
  for (let i = 0; i < points.length - 1; i++) {
    const p0 = points[i === 0 ? i : i - 1]
    const p1 = points[i]
    const p2 = points[i + 1]
    const p3 = points[i + 2 < points.length ? i + 2 : i + 1]

    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6
    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6

    d += ` C ${cp1x.toFixed(1)} ${cp1y.toFixed(1)}, ${cp2x.toFixed(1)} ${cp2y.toFixed(1)}, ${p2.x.toFixed(1)} ${p2.y.toFixed(1)}`
  }
  return d
}

const liveWaves = computed(() => {
  return {
    cyan: createSplinePath(speedHistory.value, 60),
    orange: createSplinePath(rpmHistory.value, 60),
    pink: createSplinePath(gHistory.value, 60),
  }
})

// G-Force Friction Circle Calculations (scale 2.0 G)
const maxGScale = 2.0
const gBallX = computed(() => {
  const clamped = Math.min(Math.max(lateralG.value / maxGScale, -1), 1)
  return 50 + clamped * 42
})
const gBallY = computed(() => {
  const clamped = Math.min(Math.max(-longitudinalG.value / maxGScale, -1), 1)
  return 50 + clamped * 42
})

// Drift Ball Calculations (-45° to +45° scale)
const driftBallPercent = computed(() => {
  const angle = driftAngleDeg.value
  const normalized = Math.min(Math.max(angle / 45, -1), 1)
  return 50 + normalized * 42
})

// Fuel percentage
const fuelPercent = computed(() => Math.min(Math.max((telemetry.value.fuel ?? 1) * 100, 0), 100))
</script>

<template>
  <div class="glass-panel rounded-3xl p-6 flex flex-col justify-between shadow-2xl relative overflow-hidden h-full">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-indigo-900/40 pb-3">
      <div>
        <div class="text-[11px] font-mono tracking-widest text-sky-400 font-bold uppercase">STATISTICS & DYNAMICS</div>
        <div class="text-[9px] font-mono text-neutral-400 uppercase">WAVE HARMONICS · G-FORCE & DRIFT RADAR</div>
      </div>
      <div class="text-right">
        <div class="text-xs font-mono font-bold text-neutral-300">
          #{{ telemetry.lapNumber || 1 }} LAP
        </div>
        <div class="text-[9px] font-mono text-neutral-500">P{{ telemetry.racePosition || '-' }} POSITION</div>
      </div>
    </div>

    <!-- 1. Real-Time Triple Wave Harmonic Graph -->
    <div class="my-2.5">
      <div class="flex justify-between text-[10px] font-mono text-neutral-400 font-bold mb-1">
        <span>DYNAMIC TRAJECTORY WAVES</span>
        <div class="flex items-center space-x-2 text-[9px] font-bold">
          <span class="text-cyan-400">● SPEED</span>
          <span class="text-amber-400">● RPM</span>
          <span class="text-pink-400">● G-LOAD</span>
        </div>
      </div>
      <div class="h-24 bg-slate-950/70 rounded-2xl border border-indigo-950 p-2 flex items-center justify-center overflow-hidden relative shadow-inner">
        <svg class="w-full h-full" viewBox="0 0 300 70" preserveAspectRatio="none">
          <line x1="0" y1="15" x2="300" y2="15" stroke="#1e2243" stroke-width="0.75" stroke-dasharray="3,3" />
          <line x1="0" y1="35" x2="300" y2="35" stroke="#1e2243" stroke-width="0.75" stroke-dasharray="3,3" />
          <line x1="0" y1="55" x2="300" y2="55" stroke="#1e2243" stroke-width="0.75" stroke-dasharray="3,3" />

          <!-- Cyan Wave (Speed) -->
          <path :d="liveWaves.cyan" fill="transparent" stroke="#00d4ff" stroke-width="2.5" class="drop-shadow-[0_0_8px_rgba(0,212,255,0.8)]" />
          <!-- Orange Wave (RPM) -->
          <path :d="liveWaves.orange" fill="transparent" stroke="#f59e0b" stroke-width="2" class="drop-shadow-[0_0_8px_rgba(245,158,11,0.7)]" />
          <!-- Pink Wave (G-Force Load) -->
          <path :d="liveWaves.pink" fill="transparent" stroke="#ec4899" stroke-width="2" class="drop-shadow-[0_0_8px_rgba(236,72,153,0.7)]" />
        </svg>
      </div>
    </div>

    <!-- 2. Integrated G-Force Friction Radar & Drift Horizon -->
    <div class="grid grid-cols-12 gap-3 my-2 items-center bg-slate-950/60 rounded-2xl p-3 border border-indigo-950">
      <!-- G-Force Radar (5 Cols) -->
      <div class="col-span-5 flex flex-col items-center">
        <div class="text-[9px] font-mono font-bold text-sky-300 mb-1">G-FORCE RADAR</div>
        <div class="relative w-28 h-28 bg-slate-950 rounded-full border border-indigo-900/60 flex items-center justify-center shadow-inner overflow-hidden">
          <div class="absolute w-22 h-22 rounded-full border border-indigo-900/40" />
          <div class="absolute w-14 h-14 rounded-full border border-indigo-800/60" />
          <div class="absolute inset-x-0 h-px bg-cyan-500/40" />
          <div class="absolute inset-y-0 w-px bg-cyan-500/40" />

          <!-- G-Ball -->
          <div
            class="absolute w-3.5 h-3.5 rounded-full transition-all duration-75 transform -translate-x-1/2 -translate-y-1/2 flex items-center justify-center shadow-md"
            :class="totalG > 1.2 ? 'bg-rose-500 shadow-[0_0_10px_rgba(244,63,94,1)]' : 'bg-cyan-400 shadow-[0_0_8px_rgba(56,189,248,1)]'"
            :style="{ left: `${gBallX}%`, top: `${gBallY}%` }"
          >
            <div class="w-1 h-1 rounded-full bg-white" />
          </div>
          <div class="w-1.5 h-1.5 rounded-full bg-indigo-400 pointer-events-none" />
        </div>
        <div class="text-[9px] font-mono text-cyan-300 font-bold mt-1 tabular-nums">
          {{ totalG.toFixed(2) }} G TOTAL
        </div>
      </div>

      <!-- Drift Ball Horizon (7 Cols) -->
      <div class="col-span-7 flex flex-col justify-between h-full pl-1">
        <div class="flex justify-between items-center text-[9px] font-mono">
          <span class="text-sky-300 font-bold">CHASSIS SLIP ANGLE</span>
          <span
            class="font-black tabular-nums"
            :class="isDrifting ? 'text-pink-400 animate-pulse' : 'text-neutral-300'"
          >
            {{ Math.abs(driftAngleDeg) > 0.5 ? `${driftAngleDeg.toFixed(1)}°` : '0.0°' }}
          </span>
        </div>

        <div class="relative w-full h-16 bg-slate-950/90 rounded-xl border border-indigo-900/80 p-2 flex flex-col justify-between my-1">
          <div class="text-[8px] font-mono text-neutral-400 flex justify-between">
            <span>TRAJECTORY</span>
            <span :class="isDrifting ? 'text-pink-400 font-bold' : 'text-neutral-500'">
              {{ driftAngleDeg > 3 ? 'OVERSTEER R' : driftAngleDeg < -3 ? 'OVERSTEER L' : 'NEUTRAL' }}
            </span>
          </div>

          <!-- Drift Ball Track -->
          <div class="relative w-full h-5 bg-slate-900 rounded-full border border-indigo-900 flex items-center px-1 overflow-hidden">
            <div class="absolute left-1/2 -translate-x-1/2 w-0.5 h-full bg-indigo-500 z-10" />
            <!-- Trail -->
            <div
              v-if="driftAngleDeg < 0"
              class="absolute right-1/2 h-full bg-gradient-to-l from-indigo-500 to-rose-500 rounded-l-full opacity-80"
              :style="{ width: `${Math.abs(driftAngleDeg) * 1.1}%` }"
            />
            <div
              v-if="driftAngleDeg > 0"
              class="absolute left-1/2 h-full bg-gradient-to-r from-indigo-500 to-rose-500 rounded-r-full opacity-80"
              :style="{ width: `${driftAngleDeg * 1.1}%` }"
            />
            <!-- Puck -->
            <div
              class="absolute w-4 h-4 rounded-full border-2 border-white transform -translate-x-1/2 z-20 flex items-center justify-center transition-all duration-75"
              :class="isDrifting ? 'bg-pink-500 shadow-[0_0_12px_rgba(244,63,94,1)]' : 'bg-cyan-400 shadow-[0_0_8px_rgba(56,189,248,1)]'"
              :style="{ left: `${driftBallPercent}%` }"
            >
              <div class="w-1 h-1 rounded-full bg-white" />
            </div>
          </div>

          <div class="flex justify-between text-[7px] font-mono text-neutral-500 px-1">
            <span>-45°</span>
            <span>0°</span>
            <span>+45°</span>
          </div>
        </div>

        <div class="text-[8px] font-mono text-neutral-400">
          LAT: <span class="text-cyan-300 font-bold">{{ lateralG >= 0 ? `+${lateralG.toFixed(2)}` : lateralG.toFixed(2) }}G</span>
          · LON: <span class="text-pink-400 font-bold">{{ longitudinalG >= 0 ? `+${longitudinalG.toFixed(2)}` : longitudinalG.toFixed(2) }}G</span>
        </div>
      </div>
    </div>

    <!-- Fuel Tank Reserve Gauge -->
    <div class="pt-3 border-t border-indigo-900/40">
      <div class="flex justify-between items-center mb-1.5 text-[10px] font-mono font-bold">
        <span class="text-neutral-400 tracking-wider">FUEL RESERVE CAPACITY</span>
        <span class="text-emerald-400 font-bold">{{ Math.round(fuelPercent) }}%</span>
      </div>
      <div class="w-full h-2.5 bg-slate-950 rounded-full border border-indigo-950 overflow-hidden p-0.5">
        <div
          class="h-full rounded-full transition-all duration-150"
          :class="fuelPercent < 15 ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.8)]' : 'bg-gradient-to-r from-teal-500 to-emerald-400'"
          :style="{ width: `${fuelPercent}%` }"
        />
      </div>
    </div>
  </div>
</template>
