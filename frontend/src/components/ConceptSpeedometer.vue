<script setup lang="ts">
import { computed } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { telemetry, speedKmh, gearDisplay } = useTelemetryStore()

// Ambient temperature in Celsius (returns 0°C when idle or reading is 0)
const ambientTempC = computed(() => {
  const rawF = telemetry.value.tireTemp?.[0] ?? 0
  return rawF > 0 ? Math.round((rawF - 32) * (5 / 9)) : 0
})

// Boost in Bar (0 to 2.5 bar typical scale)
const boostBar = computed(() => {
  const b = (telemetry.value.boost || 0) * 0.0689476
  return Math.max(b, 0)
})

const boostPercent = computed(() => {
  return Math.min(Math.max((boostBar.value / 2.2) * 100, 0), 100)
})

const powerKw = computed(() => {
  return Math.max(Math.round((telemetry.value.power || 0) / 1000), 0)
})

const torqueNm = computed(() => {
  return Math.max(Math.round(telemetry.value.torque || 0), 0)
})

// Max ratings for gradient percentage bars
const powerPercent = computed(() => {
  return Math.min(Math.max((powerKw.value / 650) * 100, 0), 100)
})

const torquePercent = computed(() => {
  return Math.min(Math.max((torqueNm.value / 950) * 100, 0), 100)
})

// Circular Tachometer Arc Calculation (240 deg sweep, mouth at bottom)
const dialRadius = 116
const circumference = 2 * Math.PI * dialRadius
const arcLength = circumference * (240 / 360)
const gapLength = circumference * (120 / 360)

// Dynamic RPM Percentage (0..1)
const rpmNorm = computed(() => {
  const currentRpm = telemetry.value.currentEngineRpm || 0
  const maxRpm = telemetry.value.engineMaxRpm || 8000
  return Math.min(Math.max(currentRpm / maxRpm, 0), 1)
})

// Active stroke length: sweeps dynamically with engine RPM
const activeDashArray = computed(() => {
  const activeLength = rpmNorm.value * arcLength
  return `${activeLength} ${circumference}`
})

// Near redline check (> 88% of max RPM)
const isRedline = computed(() => rpmNorm.value > 0.88)
</script>

<template>
  <div class="glass-panel rounded-3xl p-6 flex flex-col justify-between relative overflow-hidden shadow-2xl min-h-[580px]">
    <!-- Ambient cyan glow behind dial -->
    <div class="absolute -top-16 -left-16 w-64 h-64 bg-cyan-500/10 rounded-full blur-3xl pointer-events-none" />

    <!-- Top Speedometer Title Bar -->
    <div class="flex items-center justify-between border-b border-indigo-900/40 pb-3">
      <div>
        <div class="text-[11px] font-mono tracking-widest text-sky-400 font-bold uppercase">SPEEDOMETER</div>
        <div class="text-[9px] font-mono text-neutral-400 uppercase">TELEMETRY V0.1 · METRIC DASH</div>
      </div>
      <div class="text-right">
        <div class="text-xs font-mono font-bold text-neutral-300">
          {{ ambientTempC }}°C
        </div>
        <div class="text-[9px] font-mono text-neutral-500">AMBIENT</div>
      </div>
    </div>

    <!-- 1. Swept Tachometer Dial with Centered km/h & Gear -->
    <div class="relative w-full aspect-square max-w-[280px] mx-auto flex items-center justify-center my-auto">
      <svg class="w-full h-full" viewBox="0 0 280 280">
        <!-- Outer tick track ring (rotated 150deg so the 120deg mouth is centered at bottom) -->
        <circle
          cx="140"
          cy="140"
          :r="dialRadius"
          fill="transparent"
          stroke="#131735"
          stroke-width="12"
          stroke-linecap="round"
          :stroke-dasharray="`${arcLength} ${gapLength}`"
          transform="rotate(150 140 140)"
        />
        <!-- Active Tachometer RPM Glow Arc (sweeps with engine RPM, shifts to rose near redline) -->
        <circle
          cx="140"
          cy="140"
          :r="dialRadius"
          fill="transparent"
          :stroke="isRedline ? '#f43f5e' : '#00d4ff'"
          stroke-width="12"
          stroke-linecap="round"
          class="transition-all duration-75"
          :class="isRedline ? 'drop-shadow-[0_0_16px_rgba(244,63,94,0.9)]' : 'drop-shadow-[0_0_14px_rgba(0,212,255,0.85)]'"
          :stroke-dasharray="activeDashArray"
          transform="rotate(150 140 140)"
        />
      </svg>

      <!-- Speed & Gear Digital Hub Inside the Dial -->
      <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none select-none">
        <!-- Gear Pill -->
        <div class="flex items-center space-x-1.5 px-3 py-0.5 rounded-full bg-slate-900/90 border border-indigo-700/60 shadow-md mb-0.5">
          <span class="text-[9px] font-mono text-neutral-400 font-semibold uppercase">GEAR</span>
          <span class="text-lg font-black font-mono text-amber-400 leading-none">{{ gearDisplay }}</span>
        </div>

        <!-- Big Speed Digital Value -->
        <div class="text-6xl font-black font-mono tracking-tighter text-white tabular-nums drop-shadow-md text-glow-cyan leading-none my-1">
          {{ Math.round(speedKmh) }}
        </div>

        <!-- Unit Label -->
        <span class="text-xs font-mono tracking-widest text-sky-300 font-bold uppercase">KM / H</span>

        <!-- Transmission PRND Selector -->
        <div class="flex items-center space-x-2 font-mono text-[10px] font-bold mt-2">
          <span :class="gearDisplay === 'R' ? 'text-rose-400 font-black scale-110' : 'text-slate-600'">R</span>
          <span :class="gearDisplay === 'N' ? 'text-amber-400 font-black scale-110' : 'text-slate-600'">N</span>
          <span :class="gearDisplay !== 'R' && gearDisplay !== 'N' ? 'text-sky-400 font-black scale-110' : 'text-slate-600'">D</span>
        </div>
      </div>
    </div>

    <!-- Bottom Metrics Group: Clean spacing and harmonious structure -->
    <div class="space-y-4 pt-4 border-t border-indigo-900/40">
      <!-- 2. Power (k-W) and Torque (Nm) Horizontal Bars -->
      <div class="space-y-2.5">
        <!-- k-W Bar -->
        <div class="flex items-center space-x-3">
          <span class="w-10 text-xs font-mono font-bold text-sky-300">k-W</span>
          <div class="flex-1 h-3 bg-slate-950/80 rounded-full border border-indigo-900/40 overflow-hidden p-0.5 shadow-inner">
            <div
              class="h-full rounded-full bg-gradient-to-r from-pink-500 to-rose-500 transition-all duration-75 shadow-[0_0_8px_rgba(244,63,94,0.7)]"
              :style="{ width: `${powerPercent}%` }"
            />
          </div>
          <span class="w-12 text-right font-mono font-bold text-xs text-rose-400 tabular-nums">{{ powerKw }}</span>
        </div>

        <!-- Nm Bar -->
        <div class="flex items-center space-x-3">
          <span class="w-10 text-xs font-mono font-bold text-sky-300">Nm</span>
          <div class="flex-1 h-3 bg-slate-950/80 rounded-full border border-indigo-900/40 overflow-hidden p-0.5 shadow-inner">
            <div
              class="h-full rounded-full bg-gradient-to-r from-cyan-500 to-blue-500 transition-all duration-75 shadow-[0_0_8px_rgba(6,182,212,0.7)]"
              :style="{ width: `${torquePercent}%` }"
            />
          </div>
          <span class="w-12 text-right font-mono font-bold text-xs text-cyan-400 tabular-nums">{{ torqueNm }}</span>
        </div>
      </div>

      <!-- 3. Rainbow Turbo Boost Slider -->
      <div class="space-y-1.5 pt-2">
        <div class="flex justify-between text-[10px] font-mono font-bold">
          <span class="text-sky-400 italic">min</span>
          <span class="text-sky-300 tracking-widest font-black uppercase">TURBO BOOST</span>
          <span class="text-pink-400 italic">max</span>
        </div>
        <div class="relative w-full h-3 bg-slate-950/80 rounded-full p-0.5 border border-indigo-800/40 shadow-inner overflow-hidden">
          <div
            class="h-full rounded-full bg-gradient-to-r from-cyan-400 via-indigo-400 via-pink-500 to-rose-500 transition-all duration-75 shadow-[0_0_12px_rgba(236,72,153,0.7)]"
            :style="{ width: `${boostPercent}%` }"
          />
        </div>
        <div class="flex justify-between text-[9px] font-mono text-neutral-400 px-1">
          <span>0.00 BAR</span>
          <span class="text-pink-300 font-bold">{{ boostBar.toFixed(2) }} BAR</span>
          <span>2.20 BAR</span>
        </div>
      </div>
    </div>
  </div>
</template>
