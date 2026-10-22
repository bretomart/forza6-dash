<script setup lang="ts">
import { ref, computed } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { telemetry, speedKmh, speedMph, gearDisplay, rpmPercentage } = useTelemetryStore()
const unit = ref<'kmh' | 'mph'>('kmh')

const toggleUnit = () => {
  unit.value = unit.value === 'kmh' ? 'mph' : 'kmh'
}

// 24-segment high-density LED shift lights
const totalLeds = 24

// Redline / shift blink trigger (> 92% RPM)
const isRedlining = computed(() => rpmPercentage.value >= 92)
const isShiftFlashing = computed(() => rpmPercentage.value >= 96)

const currentSpeed = computed(() => {
  return Math.round(unit.value === 'kmh' ? speedKmh.value : speedMph.value)
})

// Calculate curved arc offset for sports gauge dial
const arcRadius = 110
const circumference = 2 * Math.PI * arcRadius
// 240 degree sweep arc (-210 to +30 degrees)
const arcLength = circumference * (240 / 360)
const strokeDashoffset = computed(() => {
  const pct = Math.min(Math.max(rpmPercentage.value / 100, 0), 1)
  return arcLength - pct * arcLength
})
</script>

<template>
  <div class="relative bg-carbon border border-neutral-800/90 rounded-2xl p-5 md:p-7 bevel-cockpit flex flex-col justify-between overflow-hidden shadow-2xl">
    <!-- Redline flash overlay -->
    <div
      v-if="isShiftFlashing"
      class="absolute inset-0 bg-red-600/15 pointer-events-none animate-pulse z-20"
    />

    <!-- Header / Mode Bar -->
    <div class="flex items-center justify-between pb-3 border-b border-neutral-800/80 z-10">
      <div class="flex items-center space-x-3">
        <div class="flex items-center space-x-1.5 px-2.5 py-1 rounded bg-neutral-900 border border-neutral-700/80 text-[10px] font-mono tracking-widest text-amber-400 uppercase font-black">
          <span class="w-2 h-2 rounded-full bg-amber-400 animate-ping" v-if="telemetry.isRaceOn" />
          <span>{{ telemetry.isRaceOn ? 'RACE MODE' : 'STANDBY' }}</span>
        </div>
        <span class="text-xs font-mono font-bold tracking-widest uppercase text-neutral-400">DIGITAL COCKPIT CLUSTER</span>
      </div>

      <div class="flex items-center space-x-2">
        <button
          @click="toggleUnit"
          class="text-[11px] px-3 py-1 rounded font-mono font-bold border border-neutral-700 bg-neutral-800 hover:bg-neutral-700 text-cyan-300 transition-colors cursor-pointer tracking-wider"
          title="Toggle Unit"
        >
          {{ unit.toUpperCase() }}
        </button>
      </div>
    </div>

    <!-- F1 / GT3 Style Shift Light LED Row -->
    <div class="my-3 z-10">
      <div class="grid grid-cols-24 gap-1 p-1.5 bg-black/80 rounded-lg border border-neutral-800 shadow-inner">
        <div
          v-for="i in totalLeds"
          :key="i"
          class="h-4 rounded-xs transition-all duration-75"
          :class="[
            (i / totalLeds) * 100 <= rpmPercentage
              ? i > 20
                ? 'bg-blue-400 shadow-[0_0_10px_rgba(96,165,250,1)]'
                : i > 15
                ? 'bg-red-500 shadow-[0_0_10px_rgba(239,68,68,0.9)]'
                : i > 8
                ? 'bg-amber-400 shadow-[0_0_8px_rgba(251,191,36,0.9)]'
                : 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.9)]'
              : 'bg-neutral-900 border border-neutral-800/40'
          ]"
        />
      </div>
    </div>

    <!-- Central Sports Dial & Speed Display -->
    <div class="relative flex flex-col md:flex-row items-center justify-around py-4 z-10 gap-6">
      <!-- Tachometer Circular Arc Gauge -->
      <div class="relative w-64 h-64 flex items-center justify-center shrink-0">
        <svg class="w-full h-full -rotate-120 transform" viewBox="0 0 260 260">
          <!-- Background track -->
          <circle
            cx="130"
            cy="130"
            r="110"
            fill="transparent"
            stroke="#1c1f24"
            stroke-width="14"
            stroke-linecap="round"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="circumference - arcLength"
          />
          <!-- Active RPM sweep arc -->
          <circle
            cx="130"
            cy="130"
            r="110"
            fill="transparent"
            :stroke="isRedlining ? '#ef4444' : rpmPercentage > 75 ? '#f59e0b' : '#06b6d4'"
            stroke-width="14"
            stroke-linecap="round"
            class="transition-all duration-75"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="strokeDashoffset"
          />
        </svg>

        <!-- Center Tachometer Hub with Gear Indicator -->
        <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
          <span class="text-[10px] font-mono font-bold tracking-widest text-neutral-400 uppercase">GEAR</span>
          <div
            class="w-20 h-20 rounded-2xl border-2 flex items-center justify-center font-mono font-black text-5xl shadow-2xl transition-all duration-100"
            :class="[
              gearDisplay === 'R'
                ? 'border-red-500 bg-red-950/40 text-red-400 shadow-[0_0_20px_rgba(239,68,68,0.5)] text-glow-red'
                : gearDisplay === 'N'
                ? 'border-neutral-600 bg-neutral-900/80 text-neutral-300'
                : 'border-cyan-400/80 bg-neutral-950/80 text-cyan-300 shadow-[0_0_20px_rgba(6,182,212,0.4)] text-glow-cyan'
            ]"
          >
            {{ gearDisplay }}
          </div>
          <span class="text-[11px] font-mono text-neutral-400 font-bold mt-1 tracking-wider">
            {{ Math.round(telemetry.currentEngineRpm) }} <span class="text-neutral-500 text-[9px]">RPM</span>
          </span>
        </div>
      </div>

      <!-- Speedometer Pod -->
      <div class="flex flex-col items-center md:items-start justify-center">
        <div class="text-[11px] font-mono tracking-widest uppercase font-bold text-neutral-400">GROUND SPEED</div>
        <div class="flex items-baseline space-x-2 my-1">
          <span class="text-7xl md:text-8xl font-black font-mono tracking-tighter text-white tabular-nums drop-shadow-lg text-glow-cyan">
            {{ currentSpeed }}
          </span>
          <span class="text-xl font-black font-mono text-cyan-400 uppercase tracking-widest">
            {{ unit }}
          </span>
        </div>

        <!-- Secondary telemetry readout badges -->
        <div class="grid grid-cols-2 gap-2 mt-2 w-full">
          <div class="bg-black/60 border border-neutral-800 rounded px-2.5 py-1.5 font-mono">
            <div class="text-[9px] text-neutral-500 uppercase">MAX RPM</div>
            <div class="text-xs font-bold text-neutral-300">{{ Math.round(telemetry.engineMaxRpm) }}</div>
          </div>
          <div class="bg-black/60 border border-neutral-800 rounded px-2.5 py-1.5 font-mono">
            <div class="text-[9px] text-neutral-500 uppercase">IDLE RPM</div>
            <div class="text-xs font-bold text-neutral-300">{{ Math.round(telemetry.engineIdleRpm) }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Telemetry Ticker -->
    <div class="flex items-center justify-between pt-3 border-t border-neutral-800/80 text-xs font-mono text-neutral-400 z-10">
      <div class="flex items-center space-x-2">
        <span class="w-2 h-2 rounded-full" :class="telemetry.isRaceOn ? 'bg-emerald-400 animate-pulse' : 'bg-neutral-600'" />
        <span>STATUS: {{ telemetry.isRaceOn ? 'LIVE TELEMETRY' : 'IDLE' }}</span>
      </div>
      <div>
        <span>CAR PI: </span>
        <span class="text-amber-400 font-bold tracking-wider">{{ telemetry.carPerformanceIndex || '---' }}</span>
      </div>
    </div>
  </div>
</template>
