<script setup lang="ts">
import { useTelemetryStore } from '../stores/telemetryStore'

const {
  telemetry,
  accelPercentage,
  brakePercentage,
  clutchPercentage,
  handbrakePercentage,
  steerPercentage,
} = useTelemetryStore()
</script>

<template>
  <div class="bg-carbon border border-neutral-800/90 rounded-2xl p-5 md:p-7 bevel-cockpit flex flex-col justify-between shadow-2xl">
    <!-- Header -->
    <div class="flex items-center justify-between pb-3 border-b border-neutral-800/80">
      <div class="flex items-center space-x-2">
        <div class="w-2.5 h-2.5 rounded-full bg-indigo-500 shadow-[0_0_8px_rgba(99,102,241,0.8)]" />
        <span class="text-xs font-mono font-bold tracking-widest uppercase text-neutral-300">PEDALS & CHASSIS INPUT</span>
      </div>
      <div class="flex items-center space-x-2">
        <span
          class="px-2.5 py-0.5 rounded font-mono font-black text-[10px] tracking-wider transition-all duration-100"
          :class="handbrakePercentage > 0 ? 'bg-red-600 text-white shadow-[0_0_12px_rgba(239,68,68,0.7)] animate-pulse' : 'bg-neutral-900 border border-neutral-800 text-neutral-600'"
        >
          PARK / E-BRAKE
        </span>
      </div>
    </div>

    <!-- Sports Car Vertical Pedal Box -->
    <div class="grid grid-cols-3 gap-6 py-5 my-auto">
      <!-- Clutch -->
      <div class="flex flex-col items-center">
        <div class="text-[10px] font-mono font-bold tracking-wider text-neutral-400 mb-2 uppercase">CLUTCH</div>
        <div class="relative w-14 h-48 bg-neutral-950 rounded-xl border border-neutral-800/90 p-1.5 flex flex-col justify-end overflow-hidden shadow-inner">
          <!-- Scale markers -->
          <div class="absolute inset-y-2 left-1.5 flex flex-col justify-between pointer-events-none opacity-20">
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
          </div>
          <div
            class="w-full bg-gradient-to-t from-blue-700 to-blue-400 rounded-lg transition-all duration-75 shadow-[0_0_14px_rgba(59,130,246,0.7)]"
            :style="{ height: `${clutchPercentage}%` }"
          />
        </div>
        <div class="mt-2.5 font-mono text-sm font-black text-blue-400 tabular-nums">
          {{ Math.round(clutchPercentage) }}%
        </div>
      </div>

      <!-- Brake -->
      <div class="flex flex-col items-center">
        <div class="text-[10px] font-mono font-bold tracking-wider text-neutral-400 mb-2 uppercase">BRAKE</div>
        <div class="relative w-14 h-48 bg-neutral-950 rounded-xl border border-neutral-800/90 p-1.5 flex flex-col justify-end overflow-hidden shadow-inner">
          <!-- Scale markers -->
          <div class="absolute inset-y-2 left-1.5 flex flex-col justify-between pointer-events-none opacity-20">
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
          </div>
          <div
            class="w-full bg-gradient-to-t from-red-700 to-red-400 rounded-lg transition-all duration-75 shadow-[0_0_14px_rgba(239,68,68,0.8)]"
            :style="{ height: `${brakePercentage}%` }"
          />
        </div>
        <div class="mt-2.5 font-mono text-sm font-black text-red-400 tabular-nums">
          {{ Math.round(brakePercentage) }}%
        </div>
      </div>

      <!-- Throttle / Accel -->
      <div class="flex flex-col items-center">
        <div class="text-[10px] font-mono font-bold tracking-wider text-neutral-400 mb-2 uppercase">THROTTLE</div>
        <div class="relative w-14 h-48 bg-neutral-950 rounded-xl border border-neutral-800/90 p-1.5 flex flex-col justify-end overflow-hidden shadow-inner">
          <!-- Scale markers -->
          <div class="absolute inset-y-2 left-1.5 flex flex-col justify-between pointer-events-none opacity-20">
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
            <span class="w-2 h-0.5 bg-neutral-400"></span>
          </div>
          <div
            class="w-full bg-gradient-to-t from-emerald-700 to-emerald-400 rounded-lg transition-all duration-75 shadow-[0_0_14px_rgba(16,185,129,0.8)]"
            :style="{ height: `${accelPercentage}%` }"
          />
        </div>
        <div class="mt-2.5 font-mono text-sm font-black text-emerald-400 tabular-nums">
          {{ Math.round(accelPercentage) }}%
        </div>
      </div>
    </div>

    <!-- Steering Rack Dial -->
    <div class="pt-4 border-t border-neutral-800/80">
      <div class="flex items-center justify-between text-xs font-mono text-neutral-400 mb-2">
        <span class="tracking-widest uppercase font-bold text-neutral-300">STEERING RACK ANGLE</span>
        <span class="font-bold font-mono text-neutral-200 tabular-nums">
          {{ telemetry.steer > 0 ? `+${telemetry.steer}` : telemetry.steer }}°
          <span class="text-neutral-500 font-normal">({{ Math.round(steerPercentage) }}%)</span>
        </span>
      </div>

      <!-- Cockpit Bipolar Steering Bar with Center Notch -->
      <div class="relative w-full h-4 bg-neutral-950 rounded-full border border-neutral-800 overflow-hidden flex items-center shadow-inner">
        <div class="absolute left-1/2 -translate-x-1/2 w-1 h-full bg-neutral-500 z-10" />

        <!-- Left Steer -->
        <div
          v-if="steerPercentage < 0"
          class="absolute right-1/2 h-full bg-gradient-to-l from-indigo-500 to-cyan-400 rounded-l-full shadow-[0_0_12px_rgba(99,102,241,0.8)] transition-all duration-75"
          :style="{ width: `${Math.abs(steerPercentage) / 2}%` }"
        />

        <!-- Right Steer -->
        <div
          v-if="steerPercentage > 0"
          class="absolute left-1/2 h-full bg-gradient-to-r from-indigo-500 to-cyan-400 rounded-r-full shadow-[0_0_12px_rgba(99,102,241,0.8)] transition-all duration-75"
          :style="{ width: `${steerPercentage / 2}%` }"
        />
      </div>

      <div class="flex justify-between text-[10px] font-mono text-neutral-500 mt-1.5 font-semibold">
        <span>◄ FULL LOCK L</span>
        <span>0° CENTER</span>
        <span>FULL LOCK R ►</span>
      </div>
    </div>
  </div>
</template>
