<script setup lang="ts">
import { computed } from 'vue'
import { useTelemetryStore } from '../stores/telemetryStore'
import CarWireframe from './CarWireframe.vue'

const { telemetry, lateralG, longitudinalG, driftAngleDeg, isDrifting } = useTelemetryStore()

// Tire diagnostic values in Celsius
// When idle/no reading (packet is 0), show 0°C instead of converting 0°F to -18°C
const getTireStatus = (idx: number) => {
  const rawF = telemetry.value.tireTemp?.[idx] ?? 0
  const tempC = rawF > 0 ? Math.round((rawF - 32) * (5 / 9)) : 0
  const slip = Math.abs(telemetry.value.tireSlipRatio?.[idx] ?? 0)
  const isSlipping = slip > 0.18
  return { tempC, isSlipping, slipPct: Math.round(slip * 100) }
}

// Suspension travel percentages (0..100)
const suspTravel = computed(() => [
  Math.min(Math.max((telemetry.value.normalizedSuspensionTravel?.[0] ?? 0) * 100, 0), 100),
  Math.min(Math.max((telemetry.value.normalizedSuspensionTravel?.[1] ?? 0) * 100, 0), 100),
  Math.min(Math.max((telemetry.value.normalizedSuspensionTravel?.[2] ?? 0) * 100, 0), 100),
  Math.min(Math.max((telemetry.value.normalizedSuspensionTravel?.[3] ?? 0) * 100, 0), 100)
])
</script>

<template>
  <div class="glass-panel rounded-3xl p-6 flex flex-col justify-between shadow-2xl relative overflow-hidden min-h-[580px]">
    <!-- Perspective Cyber Grid in background -->
    <div class="absolute inset-0 perspective-grid opacity-30 pointer-events-none" />

    <!-- Ambient Neon Aura around car -->
    <div
      class="absolute inset-0 flex items-center justify-center pointer-events-none transition-colors duration-300"
    >
      <div
        class="w-72 h-96 rounded-full blur-3xl opacity-35"
        :class="isDrifting ? 'bg-pink-500' : 'bg-cyan-500'"
      />
    </div>

    <!-- Header Status -->
    <div class="flex items-center justify-between border-b border-indigo-900/40 pb-3 relative z-10">
      <div>
        <div class="text-[11px] font-mono tracking-widest text-sky-400 font-bold uppercase">CHASSIS TWIN</div>
        <div class="text-[9px] font-mono text-neutral-400 uppercase">3D TELEMETRY MATRIX · TOP-DOWN HUD</div>
      </div>
      <div
        class="px-2.5 py-1 rounded-full text-[10px] font-mono font-bold tracking-wider"
        :class="[
          isDrifting
            ? 'bg-pink-950/80 border border-pink-500 text-pink-300 shadow-[0_0_12px_rgba(244,63,94,0.6)] animate-pulse'
            : 'bg-indigo-950/60 border border-indigo-700/60 text-cyan-300'
        ]"
      >
        {{ isDrifting ? 'OVERSTEER ACTIVE' : 'NEUTRAL CHASSIS' }}
      </div>
    </div>

    <!-- Center Top-Down Car Silhouette with Corner Tire & Suspension Pods -->
    <div class="relative flex-1 flex items-center justify-center py-4 my-2 z-10">
      <!-- Top-Down 3D Sports Car Wireframe -->
      <div class="w-48 h-96 relative flex items-center justify-center">
        <CarWireframe />
      </div>

      <!-- 4 Integrated Corner Wheel & Damper Pods -->
      <!-- Front Left (FL) -->
      <div
        class="absolute top-2 left-1 bg-slate-950/90 border rounded-2xl p-2.5 backdrop-blur-md transition-all duration-150 flex items-center gap-2.5 shadow-lg"
        :class="getTireStatus(0).isSlipping ? 'border-rose-500 shadow-[0_0_14px_rgba(244,63,94,0.6)]' : 'border-indigo-900/60'"
      >
        <div>
          <div class="flex items-center justify-between text-[9px] font-mono font-bold text-neutral-400 mb-0.5 gap-2">
            <span>FL WHEEL</span>
            <span class="text-[8px] text-cyan-400">CH_01</span>
          </div>
          <div class="text-xl font-black font-mono text-white leading-tight">{{ getTireStatus(0).tempC }}°C</div>
          <div class="text-[9px] font-mono font-bold" :class="getTireStatus(0).isSlipping ? 'text-rose-400' : 'text-slate-400'">
            {{ getTireStatus(0).slipPct }}% SLIP
          </div>
        </div>

        <!-- Vertical Suspension Gauge -->
        <div class="flex flex-col items-center pl-2 border-l border-indigo-950">
          <span class="text-[8px] font-mono text-cyan-300 font-bold mb-1">{{ Math.round(suspTravel[0]) }}%</span>
          <div class="w-1.5 h-10 bg-slate-900 rounded-full border border-indigo-950 overflow-hidden flex flex-col justify-end">
            <div
              class="w-full bg-gradient-to-t from-cyan-500 to-blue-400 rounded-full transition-all duration-75"
              :style="{ height: `${suspTravel[0]}%` }"
            />
          </div>
          <span class="text-[7px] font-mono text-neutral-500 mt-0.5 uppercase font-semibold">SUS</span>
        </div>
      </div>

      <!-- Front Right (FR) -->
      <div
        class="absolute top-2 right-1 bg-slate-950/90 border rounded-2xl p-2.5 backdrop-blur-md transition-all duration-150 flex items-center gap-2.5 text-right shadow-lg"
        :class="getTireStatus(1).isSlipping ? 'border-rose-500 shadow-[0_0_14px_rgba(244,63,94,0.6)]' : 'border-indigo-900/60'"
      >
        <!-- Vertical Suspension Gauge -->
        <div class="flex flex-col items-center pr-2 border-r border-indigo-950">
          <span class="text-[8px] font-mono text-cyan-300 font-bold mb-1">{{ Math.round(suspTravel[1]) }}%</span>
          <div class="w-1.5 h-10 bg-slate-900 rounded-full border border-indigo-950 overflow-hidden flex flex-col justify-end">
            <div
              class="w-full bg-gradient-to-t from-cyan-500 to-blue-400 rounded-full transition-all duration-75"
              :style="{ height: `${suspTravel[1]}%` }"
            />
          </div>
          <span class="text-[7px] font-mono text-neutral-500 mt-0.5 uppercase font-semibold">SUS</span>
        </div>

        <div>
          <div class="flex items-center justify-between text-[9px] font-mono font-bold text-neutral-400 mb-0.5 gap-2">
            <span class="text-[8px] text-cyan-400">CH_02</span>
            <span>FR WHEEL</span>
          </div>
          <div class="text-xl font-black font-mono text-white leading-tight">{{ getTireStatus(1).tempC }}°C</div>
          <div class="text-[9px] font-mono font-bold" :class="getTireStatus(1).isSlipping ? 'text-rose-400' : 'text-slate-400'">
            {{ getTireStatus(1).slipPct }}% SLIP
          </div>
        </div>
      </div>

      <!-- Rear Left (RL) -->
      <div
        class="absolute bottom-2 left-1 bg-slate-950/90 border rounded-2xl p-2.5 backdrop-blur-md transition-all duration-150 flex items-center gap-2.5 shadow-lg"
        :class="getTireStatus(2).isSlipping ? 'border-rose-500 shadow-[0_0_14px_rgba(244,63,94,0.6)]' : 'border-indigo-900/60'"
      >
        <div>
          <div class="flex items-center justify-between text-[9px] font-mono font-bold text-neutral-400 mb-0.5 gap-2">
            <span>RL WHEEL</span>
            <span class="text-[8px] text-pink-400">CH_03</span>
          </div>
          <div class="text-xl font-black font-mono text-white leading-tight">{{ getTireStatus(2).tempC }}°C</div>
          <div class="text-[9px] font-mono font-bold" :class="getTireStatus(2).isSlipping ? 'text-rose-400' : 'text-slate-400'">
            {{ getTireStatus(2).slipPct }}% SLIP
          </div>
        </div>

        <!-- Vertical Suspension Gauge -->
        <div class="flex flex-col items-center pl-2 border-l border-indigo-950">
          <span class="text-[8px] font-mono text-rose-300 font-bold mb-1">{{ Math.round(suspTravel[2]) }}%</span>
          <div class="w-1.5 h-10 bg-slate-900 rounded-full border border-indigo-950 overflow-hidden flex flex-col justify-end">
            <div
              class="w-full bg-gradient-to-t from-pink-500 to-rose-400 rounded-full transition-all duration-75"
              :style="{ height: `${suspTravel[2]}%` }"
            />
          </div>
          <span class="text-[7px] font-mono text-neutral-500 mt-0.5 uppercase font-semibold">SUS</span>
        </div>
      </div>

      <!-- Rear Right (RR) -->
      <div
        class="absolute bottom-2 right-1 bg-slate-950/90 border rounded-2xl p-2.5 backdrop-blur-md transition-all duration-150 flex items-center gap-2.5 text-right shadow-lg"
        :class="getTireStatus(3).isSlipping ? 'border-rose-500 shadow-[0_0_14px_rgba(244,63,94,0.6)]' : 'border-indigo-900/60'"
      >
        <!-- Vertical Suspension Gauge -->
        <div class="flex flex-col items-center pr-2 border-r border-indigo-950">
          <span class="text-[8px] font-mono text-rose-300 font-bold mb-1">{{ Math.round(suspTravel[3]) }}%</span>
          <div class="w-1.5 h-10 bg-slate-900 rounded-full border border-indigo-950 overflow-hidden flex flex-col justify-end">
            <div
              class="w-full bg-gradient-to-t from-pink-500 to-rose-400 rounded-full transition-all duration-75"
              :style="{ height: `${suspTravel[3]}%` }"
            />
          </div>
          <span class="text-[7px] font-mono text-neutral-500 mt-0.5 uppercase font-semibold">SUS</span>
        </div>

        <div>
          <div class="flex items-center justify-between text-[9px] font-mono font-bold text-neutral-400 mb-0.5 gap-2">
            <span class="text-[8px] text-pink-400">CH_04</span>
            <span>RR WHEEL</span>
          </div>
          <div class="text-xl font-black font-mono text-white leading-tight">{{ getTireStatus(3).tempC }}°C</div>
          <div class="text-[9px] font-mono font-bold" :class="getTireStatus(3).isSlipping ? 'text-rose-400' : 'text-slate-400'">
            {{ getTireStatus(3).slipPct }}% SLIP
          </div>
        </div>
      </div>
    </div>

    <!-- Bottom Chassis Dynamics Metrics -->
    <div class="grid grid-cols-3 gap-3 pt-3 border-t border-indigo-900/40 relative z-10 text-xs font-mono">
      <div class="bg-slate-950/60 border border-indigo-950 p-2.5 rounded-xl">
        <div class="text-[9px] text-neutral-500 uppercase font-bold">LATERAL G</div>
        <div class="text-base font-black text-cyan-300">{{ lateralG.toFixed(2) }} G</div>
      </div>
      <div class="bg-slate-950/60 border border-indigo-950 p-2.5 rounded-xl">
        <div class="text-[9px] text-neutral-500 uppercase font-bold">LONGITUDINAL G</div>
        <div class="text-base font-black text-pink-400">{{ longitudinalG.toFixed(2) }} G</div>
      </div>
      <div class="bg-slate-950/60 border border-indigo-950 p-2.5 rounded-xl">
        <div class="text-[9px] text-neutral-500 uppercase font-bold">DRIFT SLIP</div>
        <div class="text-base font-black text-amber-400">{{ driftAngleDeg.toFixed(1) }}°</div>
      </div>
    </div>
  </div>
</template>
