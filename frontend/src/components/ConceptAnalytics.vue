<script setup lang="ts">
import { ref, computed, onUnmounted } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { telemetry, accelPercentage, brakePercentage } = useTelemetryStore()

// Real-time rolling history buffer of tire slip and traction for the "Run of Kilometers / Traction Wave"
const maxHistory = 40
const tractionHistory = ref<number[]>(Array(maxHistory).fill(25))

// Capture live incoming slip and speed to plot the live traction wave
const currentTractionLoad = computed(() => {
  const slips = telemetry.value.tireSlipRatio || [0, 0, 0, 0]
  const avgSlip = (Math.abs(slips[0]) + Math.abs(slips[1]) + Math.abs(slips[2]) + Math.abs(slips[3])) / 4
  const speed = (telemetry.value.speed || 0) * 3.6
  // Load value between 5 (low) and 50 (high traction peak)
  const load = Math.min(Math.max(avgSlip * 120 + (speed > 10 ? 15 : 5), 5), 52)
  return load
})

// Update history every 100ms when telemetry ticks
const interval = setInterval(() => {
  tractionHistory.value.push(currentTractionLoad.value)
  if (tractionHistory.value.length > maxHistory) {
    tractionHistory.value.shift()
  }
}, 100)

onUnmounted(() => {
  clearInterval(interval)
})

// SVG Smooth Cubic Bezier Path from real-time tractionHistory
const wavePath = computed(() => {
  const data = tractionHistory.value
  const width = 300
  const height = 60
  const step = width / (data.length - 1)

  // Map each value to SVG coordinates (invert Y: 0 is top, height is bottom)
  const points = data.map((val, i) => ({
    x: i * step,
    y: height - val
  }))

  if (points.length === 0) return { path: '', fillPath: '', peakPoint: { x: 150, y: 30 } }

  // Build SVG path with Catmull-Rom or cubic spline
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

  const fillD = `${d} L ${width} ${height} L 0 ${height} Z`
  
  // Find current peak point for pulsing glow dot
  let peakIdx = 0
  let minY = 999
  points.forEach((p, idx) => {
    if (p.y < minY) {
      minY = p.y
      peakIdx = idx
    }
  })

  return { path: d, fillPath: fillD, peakPoint: points[peakIdx] }
})

// 20-bar Equalizer Visualizer driven by current engine RPM & Speed
const rpmBars = computed(() => {
  const currentRpm = telemetry.value.currentEngineRpm || 0
  const maxRpm = telemetry.value.engineMaxRpm || 8000
  const norm = Math.min(Math.max(currentRpm / maxRpm, 0), 1)

  return Array.from({ length: 20 }, (_, i) => {
    const factor = Math.sin((i / 19) * Math.PI) * 0.4 + 0.6
    const heightPct = Math.min(Math.max(norm * factor * 100 + (Math.sin(i * 1.5) * 8), 6), 100)
    return Math.round(heightPct)
  })
})
</script>

<template>
  <div class="glass-panel rounded-3xl p-6 flex flex-col justify-between shadow-2xl relative overflow-hidden min-h-[580px]">
    <!-- Header -->
    <div class="flex items-center justify-between border-b border-indigo-900/40 pb-3">
      <div>
        <div class="text-[11px] font-mono tracking-widest text-sky-400 font-bold uppercase">DASHBOARD</div>
        <div class="text-[9px] font-mono text-neutral-400 uppercase">TELEMETRY ANALYTICS · LIVE CAN-BUS</div>
      </div>
      <div class="text-right">
        <div class="text-xs font-mono font-bold text-neutral-300">
          {{ ((telemetry.distanceTraveled || 0) / 1000).toFixed(1) }} KM
        </div>
        <div class="text-[9px] font-mono text-neutral-500">SESSION DISTANCE</div>
      </div>
    </div>

    <!-- 1. Rainbow Spectrum Equalizer (RPM & Throttle Harmonics) -->
    <div class="my-auto py-2">
      <div class="flex justify-between text-[10px] font-mono text-neutral-400 font-bold mb-2">
        <span>ENGINE HARMONICS SPECTRUM</span>
        <span class="text-pink-400 font-black tabular-nums">{{ Math.round(telemetry.currentEngineRpm) }} RPM</span>
      </div>
      <div class="h-28 bg-slate-950/70 rounded-2xl border border-indigo-950/80 p-3 flex items-end justify-between gap-1 shadow-inner">
        <div
          v-for="(barH, idx) in rpmBars"
          :key="idx"
          class="w-full rounded-t-sm transition-all duration-75 shadow-sm"
          :style="{
            height: `${barH}%`,
            background: `linear-gradient(to top, #4f46e5, ${idx > 13 ? '#f43f5e' : idx > 8 ? '#d946ef' : '#38bdf8'})`
          }"
        />
      </div>
    </div>

    <!-- 2. Dual Climate / Pedal Sliders (Concept-style Target Pucks) -->
    <div class="space-y-4 my-auto py-2">
      <div class="text-[10px] font-mono tracking-widest text-neutral-400 font-bold uppercase">
        PEDAL & CHASSIS INPUT SLIDERS
      </div>

      <!-- Throttle Slider with target puck -->
      <div class="flex items-center space-x-3 text-xs font-mono">
        <span class="w-16 text-emerald-400 font-bold text-[10px] uppercase">THROTTLE</span>
        <div class="relative flex-1 h-3 bg-slate-950 rounded-full border border-indigo-950 overflow-visible">
          <div
            class="h-full bg-gradient-to-r from-teal-500 to-emerald-400 rounded-full"
            :style="{ width: `${accelPercentage}%` }"
          />
          <div
            class="absolute top-1/2 -translate-y-1/2 w-4.5 h-4.5 rounded-full bg-white border-2 border-emerald-500 shadow-[0_0_8px_rgba(16,185,129,0.9)] transform -translate-x-1/2 transition-all duration-75"
            :style="{ left: `${accelPercentage}%` }"
          />
        </div>
        <span class="w-10 text-right font-bold text-neutral-300 tabular-nums">{{ Math.round(accelPercentage) }}%</span>
      </div>

      <!-- Brake Slider with target puck -->
      <div class="flex items-center space-x-3 text-xs font-mono">
        <span class="w-16 text-rose-400 font-bold text-[10px] uppercase">BRAKE</span>
        <div class="relative flex-1 h-3 bg-slate-950 rounded-full border border-indigo-950 overflow-visible">
          <div
            class="h-full bg-gradient-to-r from-pink-600 to-rose-400 rounded-full"
            :style="{ width: `${brakePercentage}%` }"
          />
          <div
            class="absolute top-1/2 -translate-y-1/2 w-4.5 h-4.5 rounded-full bg-white border-2 border-rose-500 shadow-[0_0_8px_rgba(244,63,94,0.9)] transform -translate-x-1/2 transition-all duration-75"
            :style="{ left: `${brakePercentage}%` }"
          />
        </div>
        <span class="w-10 text-right font-bold text-neutral-300 tabular-nums">{{ Math.round(brakePercentage) }}%</span>
      </div>
    </div>

    <!-- 3. Run of Kilometers / Dynamic Live Traction Wave (Direct Concept match) -->
    <div class="my-auto py-2">
      <div class="flex justify-between text-[10px] font-mono text-neutral-400 font-bold mb-2">
        <span>RUN OF KILOMETERS / TRACTION WAVE</span>
        <span class="text-pink-400 font-bold tabular-nums">LIVE TELEMETRY LOAD</span>
      </div>
      <div class="h-28 bg-slate-950/70 rounded-2xl border border-indigo-950 p-2.5 flex items-center justify-center overflow-hidden relative shadow-inner">
        <svg class="w-full h-full" viewBox="0 0 300 60" preserveAspectRatio="none">
          <defs>
            <linearGradient id="liveWaveGrad" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#ec4899" stop-opacity="0.6" />
              <stop offset="50%" stop-color="#8b5cf6" stop-opacity="0.3" />
              <stop offset="100%" stop-color="#3b82f6" stop-opacity="0.0" />
            </linearGradient>
          </defs>

          <!-- Rolling Dynamic Area Fill -->
          <path
            :d="wavePath.fillPath"
            fill="url(#liveWaveGrad)"
          />

          <!-- Live Glowing Wave Stroke -->
          <path
            :d="wavePath.path"
            fill="transparent"
            stroke="#ec4899"
            stroke-width="2.5"
            class="drop-shadow-[0_0_8px_rgba(236,72,153,0.9)]"
          />

          <!-- Peak Traction Indicator Puck (Concept Golden Dot) -->
          <circle
            :cx="wavePath.peakPoint.x"
            :cy="wavePath.peakPoint.y"
            r="4.5"
            fill="#fbbf24"
            stroke="#ffffff"
            stroke-width="1.5"
            class="drop-shadow-[0_0_8px_rgba(251,191,36,1)] animate-ping"
          />
          <circle
            :cx="wavePath.peakPoint.x"
            :cy="wavePath.peakPoint.y"
            r="4.5"
            fill="#fbbf24"
            stroke="#ffffff"
            stroke-width="1.5"
          />
        </svg>
      </div>
      <div class="flex justify-between text-[8px] font-mono text-neutral-500 mt-2 px-1">
        <span>-4.0s</span>
        <span>-3.0s</span>
        <span>-2.0s</span>
        <span>-1.0s</span>
        <span class="text-cyan-400 font-bold">CURRENT</span>
      </div>
    </div>
  </div>
</template>
