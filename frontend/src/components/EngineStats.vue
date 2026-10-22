<script setup lang="ts">
import { computed } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { telemetry, sessionDistanceMeters } = useTelemetryStore()

// Power in kW (Metric primary) and converted to Horsepower
const powerKw = computed(() => {
  const pWatts = telemetry.value.power || 0
  return Math.max(Math.round(pWatts / 1000), 0)
})

const horsepower = computed(() => {
  const pWatts = telemetry.value.power || 0
  return Math.max(Math.round(pWatts / 745.699872), 0)
})

const torqueNm = computed(() => {
  return Math.max(Math.round(telemetry.value.torque || 0), 0)
})

const torqueLbFt = computed(() => {
  return Math.max(Math.round((telemetry.value.torque || 0) * 0.737562), 0)
})

// Turbo Boost in Bar (Metric primary: 1 psi = 0.0689476 bar)
const boostBar = computed(() => {
  const b = telemetry.value.boost || 0
  const bar = b * 0.0689476
  return bar > 0 ? bar.toFixed(2) : '0.00'
})

const boostPsi = computed(() => {
  const b = telemetry.value.boost || 0
  return b > 0 ? b.toFixed(1) : '0.0'
})

const fuelPercentage = computed(() => {
  const f = telemetry.value.fuel ?? 1
  return Math.min(Math.max(f * 100, 0), 100)
})
</script>

<template>
  <div class="bg-carbon border border-neutral-800/90 rounded-2xl p-5 md:p-7 bevel-cockpit flex flex-col justify-between shadow-2xl">
    <!-- Header -->
    <div class="flex items-center justify-between pb-3 border-b border-neutral-800/80">
      <div class="flex items-center space-x-2">
        <div class="w-2.5 h-2.5 rounded-full bg-orange-400 shadow-[0_0_8px_rgba(251,146,60,0.8)]" />
        <span class="text-xs font-mono font-bold tracking-widest uppercase text-neutral-300">POWERTRAIN & RACE TELEMETRY</span>
      </div>
      <div class="text-xs font-mono font-bold text-neutral-400">
        CONFIG: <span class="text-amber-400">{{ telemetry.numCylinders || 4 }}-CYL</span>
      </div>
    </div>

    <!-- 2x2 Performance Readout Grid (Metric Primary) -->
    <div class="grid grid-cols-2 gap-4 py-4">
      <!-- Power Gauge Card (kW primary) -->
      <div class="bg-neutral-950/80 border border-neutral-800/80 rounded-xl p-4 flex flex-col justify-between relative overflow-hidden">
        <div class="text-[10px] font-mono tracking-wider uppercase text-neutral-400 font-bold">ENGINE POWER</div>
        <div class="flex items-baseline space-x-1.5 my-1">
          <span class="text-4xl font-black font-mono tracking-tighter text-white tabular-nums drop-shadow-md">
            {{ powerKw }}
          </span>
          <span class="text-xs font-mono font-black text-orange-400 uppercase">KW</span>
        </div>
        <div class="text-[10px] font-mono text-neutral-500 font-semibold">
          {{ horsepower }} BHP EQUIVALENT
        </div>
      </div>

      <!-- Torque Gauge Card (N·m primary) -->
      <div class="bg-neutral-950/80 border border-neutral-800/80 rounded-xl p-4 flex flex-col justify-between relative overflow-hidden">
        <div class="text-[10px] font-mono tracking-wider uppercase text-neutral-400 font-bold">CRANK TORQUE</div>
        <div class="flex items-baseline space-x-1.5 my-1">
          <span class="text-4xl font-black font-mono tracking-tighter text-white tabular-nums drop-shadow-md">
            {{ torqueNm }}
          </span>
          <span class="text-xs font-mono font-black text-cyan-400 uppercase">N·M</span>
        </div>
        <div class="text-[10px] font-mono text-neutral-500 font-semibold">
          {{ torqueLbFt }} LB-FT TORQUE
        </div>
      </div>

      <!-- Turbo Boost Gauge Card (Bar primary) -->
      <div class="bg-neutral-950/80 border border-neutral-800/80 rounded-xl p-4 flex flex-col justify-between relative overflow-hidden">
        <div class="text-[10px] font-mono tracking-wider uppercase text-neutral-400 font-bold">TURBO MANIFOLD</div>
        <div class="flex items-baseline space-x-1.5 my-1">
          <span class="text-4xl font-black font-mono tracking-tighter text-purple-300 tabular-nums drop-shadow-md">
            {{ boostBar }}
          </span>
          <span class="text-xs font-mono font-black text-purple-400 uppercase">BAR</span>
        </div>
        <div class="text-[10px] font-mono text-neutral-500 font-semibold">
          {{ boostPsi }} PSI BOOST
        </div>
      </div>

      <!-- Fuel Tank Cell -->
      <div class="bg-neutral-950/80 border border-neutral-800/80 rounded-xl p-4 flex flex-col justify-between relative overflow-hidden">
        <div class="flex justify-between items-center text-[10px] font-mono tracking-wider uppercase font-bold text-neutral-400">
          <span>FUEL RESERVE</span>
          <span
            v-if="fuelPercentage < 15"
            class="text-[9px] text-red-400 animate-pulse font-black"
          >
            LOW
          </span>
        </div>
        <div class="flex items-baseline space-x-1.5 my-1">
          <span
            class="text-4xl font-black font-mono tracking-tighter tabular-nums"
            :class="fuelPercentage < 15 ? 'text-red-400 text-glow-red animate-pulse' : 'text-white'"
          >
            {{ Math.round(fuelPercentage) }}
          </span>
          <span class="text-xs font-mono font-black text-emerald-400 uppercase">%</span>
        </div>
        <!-- Fuel Level Horizontal Gauge -->
        <div class="w-full bg-neutral-900 h-2 rounded-full overflow-hidden border border-neutral-800 mt-1">
          <div
            class="h-full rounded-full transition-all duration-150"
            :class="fuelPercentage < 15 ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.8)]' : fuelPercentage < 35 ? 'bg-amber-400' : 'bg-emerald-400'"
            :style="{ width: `${fuelPercentage}%` }"
          />
        </div>
      </div>
    </div>

    <!-- Race Sector Readouts -->
    <div class="grid grid-cols-3 gap-2 pt-3 border-t border-neutral-800/80 text-xs font-mono">
      <div class="bg-black/50 p-2.5 rounded-lg border border-neutral-800/80 flex flex-col justify-between">
        <div class="text-[9px] text-neutral-500 uppercase font-bold">CURRENT LAP</div>
        <div class="text-base font-black text-neutral-200">#{{ telemetry.lapNumber || 1 }}</div>
      </div>
      <div class="bg-black/50 p-2.5 rounded-lg border border-neutral-800/80 flex flex-col justify-between">
        <div class="text-[9px] text-neutral-500 uppercase font-bold">GRID POSITION</div>
        <div class="text-base font-black text-amber-400">P{{ telemetry.racePosition || '-' }}</div>
      </div>
      <div class="bg-black/50 p-2.5 rounded-lg border border-neutral-800/80 flex flex-col justify-between">
        <div class="text-[9px] text-neutral-500 uppercase font-bold">ODOMETER</div>
        <div class="text-base font-black text-neutral-200">{{ ((sessionDistanceMeters || telemetry.distanceTraveled || 0) / 1000).toFixed(2) }} <span class="text-[10px] text-neutral-400 font-normal">KM</span></div>
      </div>
    </div>
  </div>
</template>
