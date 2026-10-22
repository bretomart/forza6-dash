<script setup lang="ts">
import { computed } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'

const { telemetry } = useTelemetryStore()

interface TireCardData {
  key: string
  label: string
  pos: string
  index: number
}

const tires: TireCardData[] = [
  { key: 'fl', label: 'FRONT LEFT', pos: 'FL', index: 0 },
  { key: 'fr', label: 'FRONT RIGHT', pos: 'FR', index: 1 },
  { key: 'rl', label: 'REAR LEFT', pos: 'RL', index: 2 },
  { key: 'rr', label: 'REAR RIGHT', pos: 'RR', index: 3 },
]

// Temperature gradient and alert colors (evaluated in Celsius)
const getTempColorC = (tempC: number) => {
  if (tempC < 70) return 'text-sky-400'
  if (tempC <= 105) return 'text-emerald-400'
  if (tempC <= 125) return 'text-amber-400'
  return 'text-red-500 text-glow-red animate-pulse'
}

const getSuspensionPercent = (index: number) => {
  const travel = telemetry.value.normalizedSuspensionTravel?.[index] ?? 0
  return Math.min(Math.max(travel * 100, 0), 100)
}

const isSlipping = (index: number) => {
  const slip = Math.abs(telemetry.value.tireSlipRatio?.[index] ?? 0)
  return slip > 0.18
}

const anySlipActive = computed(() => {
  return [0, 1, 2, 3].some((i) => isSlipping(i))
})
</script>

<template>
  <div class="relative bg-carbon border border-neutral-800/90 rounded-2xl p-5 md:p-7 bevel-cockpit flex flex-col justify-between shadow-2xl">
    <!-- Header -->
    <div class="flex items-center justify-between pb-3 border-b border-neutral-800/80">
      <div class="flex items-center space-x-2">
        <div class="w-2.5 h-2.5 rounded-full bg-amber-400 shadow-[0_0_8px_rgba(251,191,36,0.8)]" />
        <span class="text-xs font-mono font-bold tracking-widest uppercase text-neutral-300">TIRE THERMALS & CHASSIS GEOMETRY</span>
      </div>
      <div
        class="px-2.5 py-0.5 rounded font-mono font-black text-[10px] tracking-wider transition-all duration-150"
        :class="anySlipActive ? 'bg-red-600 text-white shadow-[0_0_12px_rgba(239,68,68,0.7)] animate-pulse' : 'bg-neutral-900 border border-neutral-800 text-emerald-400'"
      >
        {{ anySlipActive ? 'TRACTION ALERT: SLIP DETECTED' : 'TRACTION: OPTIMAL' }}
      </div>
    </div>

    <!-- 2x2 Car Top-Down Chassis View Layout -->
    <div class="relative py-4">
      <!-- Minimalist vehicle silhouette in background center -->
      <div class="absolute inset-0 flex items-center justify-center pointer-events-none opacity-15">
        <div class="w-24 h-48 border-2 border-dashed border-cyan-400 rounded-3xl flex flex-col justify-between items-center py-4">
          <div class="text-[9px] font-mono text-cyan-400 font-bold uppercase">FRONT</div>
          <div class="w-12 h-0.5 bg-cyan-400" />
          <div class="text-[9px] font-mono text-cyan-400 font-bold uppercase">REAR</div>
        </div>
      </div>

      <!-- 2x2 Grid of Wheels -->
      <div class="grid grid-cols-2 gap-4 relative z-10">
        <div
          v-for="tire in tires"
          :key="tire.key"
          class="relative bg-neutral-950/85 rounded-xl p-3.5 border transition-all duration-100 flex flex-col justify-between"
          :class="[
            isSlipping(tire.index)
              ? 'border-red-500 bg-red-950/30 shadow-[0_0_18px_rgba(239,68,68,0.4)]'
              : 'border-neutral-800 hover:border-neutral-700'
          ]"
        >
          <!-- Corner Position Stamp -->
          <div class="flex justify-between items-center mb-1">
            <span class="text-xs font-mono font-black text-neutral-300">{{ tire.pos }}</span>
            <span
              class="text-[10px] font-mono font-black px-1.5 py-0.5 rounded"
              :class="isSlipping(tire.index) ? 'bg-red-600 text-white animate-bounce' : 'text-neutral-500'"
            >
              {{ ((telemetry.tireSlipRatio?.[tire.index] ?? 0) * 100).toFixed(0) }}% SLIP
            </span>
          </div>

          <!-- Tire Temperature (Metric Primary: Celsius) -->
          <div class="flex items-baseline space-x-1 my-1">
            <span
              class="text-3xl font-black font-mono tracking-tight tabular-nums"
              :class="getTempColorC(Math.round(((telemetry.tireTemp?.[tire.index] ?? 0) - 32) * (5 / 9)))"
            >
              {{ Math.round(((telemetry.tireTemp?.[tire.index] ?? 0) - 32) * (5 / 9)) }}
            </span>
            <span class="text-xs font-mono text-cyan-400 font-bold">°C</span>
            <span class="text-xs font-mono text-neutral-500 ml-1.5">
              ({{ Math.round(telemetry.tireTemp?.[tire.index] ?? 0) }}°F)
            </span>
          </div>

          <!-- Suspension Damper Compression Bar -->
          <div class="mt-2">
            <div class="flex justify-between text-[9px] font-mono text-neutral-400 mb-1">
              <span>SUSPENSION TRAVEL</span>
              <span class="font-bold text-neutral-200">{{ Math.round(getSuspensionPercent(tire.index)) }}%</span>
            </div>
            <div class="w-full bg-neutral-900 h-2 rounded-full overflow-hidden border border-neutral-800">
              <div
                class="h-full rounded-full transition-all duration-75"
                :class="
                  getSuspensionPercent(tire.index) > 85
                    ? 'bg-red-500 shadow-[0_0_8px_rgba(239,68,68,0.8)]'
                    : getSuspensionPercent(tire.index) > 60
                    ? 'bg-amber-400'
                    : 'bg-emerald-400'
                "
                :style="{ width: `${getSuspensionPercent(tire.index)}%` }"
              />
            </div>
          </div>

          <!-- Wheel Angular Speed -->
          <div class="flex justify-between text-[9px] font-mono text-neutral-500 mt-2 pt-1 border-t border-neutral-900">
            <span>ROTATION</span>
            <span class="font-semibold text-neutral-400">{{ (telemetry.wheelRotationSpeed?.[tire.index] ?? 0).toFixed(1) }} rad/s</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Thermal Color Index Bar (Metric Scale) -->
    <div class="flex items-center justify-between pt-3 border-t border-neutral-800/80 text-[10px] font-mono text-neutral-500">
      <span>TEMPERATURE INDEX (°C):</span>
      <div class="flex items-center space-x-3 font-semibold">
        <span class="text-sky-400">&lt;70° COLD</span>
        <span class="text-emerald-400">70-105° OPTIMAL</span>
        <span class="text-amber-400">106-125° HOT</span>
        <span class="text-red-400">&gt;125° OVERHEAT</span>
      </div>
    </div>
  </div>
</template>
