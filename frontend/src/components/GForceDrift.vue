<script setup lang="ts">
import { computed } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { lateralG, longitudinalG, totalG, driftAngleDeg, isDrifting } = useTelemetryStore()

// Maximum G display scale (2.0 G radius)
const maxGScale = 2.0

// Map lateral (X) and longitudinal (Y) G to 2D percentage coordinates
const gBallX = computed(() => {
  const clamped = Math.min(Math.max(lateralG.value / maxGScale, -1), 1)
  return 50 + clamped * 42
})

const gBallY = computed(() => {
  const clamped = Math.min(Math.max(-longitudinalG.value / maxGScale, -1), 1)
  return 50 + clamped * 42
})

// Drift Ball position: maps drift slip angle (-45° to +45°) across an arc or horizontal band
const driftBallPercent = computed(() => {
  const angle = driftAngleDeg.value
  const normalized = Math.min(Math.max(angle / 45, -1), 1)
  return 50 + normalized * 42
})
</script>

<template>
  <div class="glass-panel rounded-3xl p-6 shadow-2xl relative overflow-hidden">
    <!-- Ambient Cyan/Magenta Glow Aura in background -->
    <div class="absolute -top-12 left-1/4 w-80 h-32 bg-cyan-500/10 rounded-full blur-3xl pointer-events-none" />
    <div class="absolute -top-12 right-1/4 w-80 h-32 bg-pink-500/10 rounded-full blur-3xl pointer-events-none" />

    <!-- Header -->
    <div class="flex items-center justify-between pb-3 border-b border-indigo-900/40 relative z-10">
      <div class="flex items-center space-x-3">
        <div class="w-3 h-3 rounded-full bg-gradient-to-r from-cyan-400 to-pink-500 shadow-[0_0_10px_rgba(56,189,248,0.8)]" />
        <div>
          <span class="text-xs font-mono font-black tracking-widest uppercase text-sky-400">
            DYNAMICS COCKPIT: G-FORCE FRICTION RADAR & DRIFT HORIZON
          </span>
          <div class="text-[9px] font-mono text-neutral-400 uppercase">
            6-AXIS IMU DATA STREAM · REALTIME VECTOR MAPPING
          </div>
        </div>
      </div>

      <div
        class="px-3 py-1 rounded-full font-mono font-black text-[10px] tracking-wider transition-all duration-150"
        :class="[
          isDrifting
            ? 'bg-pink-950/80 border border-pink-500 text-pink-300 shadow-[0_0_14px_rgba(244,63,94,0.7)] animate-pulse'
            : 'bg-indigo-950/60 border border-indigo-800 text-cyan-300'
        ]"
      >
        {{ isDrifting ? 'DRIFT ENGAGED' : 'GRIP TRAJECTORY' }}
      </div>
    </div>

    <!-- Dual Cyber HUD Instrument Row -->
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8 py-5 items-center relative z-10">
      <!-- 1. Cyber G-Force Friction Circle Radar (Concept Glass Design) -->
      <div class="flex flex-col items-center">
        <div class="flex justify-between w-full text-[10px] font-mono text-neutral-400 font-bold mb-2">
          <span class="text-sky-300">G-FORCE FRICTION RADAR</span>
          <span class="text-pink-400 font-black tabular-nums text-glow-magenta">{{ totalG.toFixed(2) }} G RESULTANT</span>
        </div>

        <div class="relative w-48 h-48 bg-slate-950/90 rounded-full border border-indigo-800/40 p-2 flex items-center justify-center shadow-inner overflow-hidden">
          <!-- Cyber Concentric Target Rings -->
          <div class="absolute w-40 h-40 rounded-full border border-indigo-900/50" />
          <div class="absolute w-28 h-28 rounded-full border border-indigo-800/60" />
          <div class="absolute w-14 h-14 rounded-full border border-cyan-500/40" />

          <!-- Neon Crosshairs -->
          <div class="absolute inset-x-0 h-px bg-gradient-to-r from-transparent via-cyan-500/60 to-transparent" />
          <div class="absolute inset-y-0 w-px bg-gradient-to-b from-transparent via-cyan-500/60 to-transparent" />

          <!-- Cardinal Labels -->
          <span class="absolute top-1 text-[8px] font-mono text-neutral-400 font-bold">BRAKE</span>
          <span class="absolute bottom-1 text-[8px] font-mono text-neutral-400 font-bold">ACCEL</span>
          <span class="absolute left-2 text-[8px] font-mono text-neutral-400 font-bold">L</span>
          <span class="absolute right-2 text-[8px] font-mono text-neutral-400 font-bold">R</span>

          <!-- High-Glow Target Orb / G-Ball -->
          <div
            class="absolute w-5 h-5 rounded-full transition-all duration-75 transform -translate-x-1/2 -translate-y-1/2 flex items-center justify-center shadow-lg"
            :class="[
              totalG > 1.2
                ? 'bg-rose-500 shadow-[0_0_16px_rgba(244,63,94,1)]'
                : totalG > 0.7
                ? 'bg-amber-400 shadow-[0_0_14px_rgba(251,191,36,0.9)]'
                : 'bg-cyan-400 shadow-[0_0_14px_rgba(56,189,248,0.9)]'
            ]"
            :style="{ left: `${gBallX}%`, top: `${gBallY}%` }"
          >
            <div class="w-1.5 h-1.5 rounded-full bg-white" />
          </div>

          <!-- Center Reference Pivot -->
          <div class="w-2 h-2 rounded-full bg-indigo-400 pointer-events-none" />
        </div>

        <div class="flex justify-between w-full text-[10px] font-mono text-neutral-400 mt-2.5">
          <span class="bg-slate-950/70 border border-indigo-950 px-2.5 py-1 rounded-md">
            LAT: <strong class="text-cyan-300">{{ lateralG >= 0 ? `+${lateralG.toFixed(2)}` : lateralG.toFixed(2) }} G</strong>
          </span>
          <span class="bg-slate-950/70 border border-indigo-950 px-2.5 py-1 rounded-md">
            LON: <strong class="text-pink-400">{{ longitudinalG >= 0 ? `+${longitudinalG.toFixed(2)}` : longitudinalG.toFixed(2) }} G</strong>
          </span>
        </div>
      </div>

      <!-- 2. Cyber Drift Horizon Track (Concept Cyber Glass styling) -->
      <div class="flex flex-col justify-between h-full">
        <div class="flex justify-between w-full text-[10px] font-mono text-neutral-400 font-bold mb-2">
          <span class="text-sky-300">CHASSIS SLIP ANGLE / DRIFT RADIAL</span>
          <span
            class="font-black tabular-nums transition-colors"
            :class="isDrifting ? 'text-pink-400 font-bold text-glow-magenta' : 'text-neutral-300'"
          >
            {{ Math.abs(driftAngleDeg) > 0.5 ? `${driftAngleDeg.toFixed(1)}°` : '0.0°' }}
          </span>
        </div>

        <div class="relative w-full h-32 bg-slate-950/80 rounded-2xl border border-indigo-800/40 p-4 flex flex-col justify-between shadow-inner overflow-hidden my-auto">
          <!-- Drift Readout and Trajectory Status -->
          <div class="flex items-baseline justify-between z-10">
            <span class="text-xs font-mono font-bold text-neutral-400">CHASSIS SLIP ANGLE</span>
            <span
              class="text-3xl font-black font-mono tracking-tight tabular-nums"
              :class="isDrifting ? 'text-pink-400 text-glow-magenta' : 'text-white'"
            >
              {{ Math.abs(Math.round(driftAngleDeg)) }}°
            </span>
          </div>

          <!-- Drift Ball Track with Neon Trails -->
          <div class="relative w-full h-7 bg-slate-900 rounded-full border border-indigo-900 flex items-center px-1 overflow-hidden">
            <!-- Center Reference Line -->
            <div class="absolute left-1/2 -translate-x-1/2 w-1 h-full bg-indigo-500 z-10" />

            <!-- Dynamic Neon Slip Trail -->
            <div
              v-if="driftAngleDeg < 0"
              class="absolute right-1/2 h-full bg-gradient-to-l from-indigo-500 via-pink-500 to-rose-500 rounded-l-full opacity-80 transition-all duration-75 shadow-[0_0_12px_rgba(244,63,94,0.8)]"
              :style="{ width: `${Math.abs(driftAngleDeg) * 1.1}%` }"
            />
            <div
              v-if="driftAngleDeg > 0"
              class="absolute left-1/2 h-full bg-gradient-to-r from-indigo-500 via-pink-500 to-rose-500 rounded-r-full opacity-80 transition-all duration-75 shadow-[0_0_12px_rgba(244,63,94,0.8)]"
              :style="{ width: `${driftAngleDeg * 1.1}%` }"
            />

            <!-- Neon Drift Ball Puck (Matching concept Target style) -->
            <div
              class="absolute w-5 h-5 rounded-full border-2 border-white transform -translate-x-1/2 transition-all duration-75 z-20 shadow-lg flex items-center justify-center"
              :class="[
                isDrifting
                  ? 'bg-pink-500 shadow-[0_0_16px_rgba(244,63,94,1)] animate-bounce'
                  : 'bg-cyan-400 shadow-[0_0_12px_rgba(56,189,248,0.9)]'
              ]"
              :style="{ left: `${driftBallPercent}%` }"
            >
              <div class="w-1.5 h-1.5 rounded-full bg-white" />
            </div>
          </div>

          <!-- Scale Markers -->
          <div class="flex justify-between text-[9px] font-mono text-neutral-400 font-semibold px-1">
            <span>-45° L</span>
            <span>-20°</span>
            <span class="text-sky-300 font-bold">0° CENTER</span>
            <span>+20°</span>
            <span>+45° R</span>
          </div>
        </div>

        <div class="flex justify-between w-full text-[10px] font-mono text-neutral-400 mt-2.5">
          <span class="bg-slate-950/70 border border-indigo-950 px-2.5 py-1 rounded-md">
            TRAJECTORY: <strong class="text-amber-400">{{ driftAngleDeg > 3 ? 'OVERSTEER RIGHT' : driftAngleDeg < -3 ? 'OVERSTEER LEFT' : 'NEUTRAL' }}</strong>
          </span>
          <span class="bg-slate-950/70 border border-indigo-950 px-2.5 py-1 rounded-md">
            SCALE: <strong class="text-cyan-300">±45° LOCK</strong>
          </span>
        </div>
      </div>
    </div>
  </div>
</template>
