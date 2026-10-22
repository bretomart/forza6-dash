<script setup lang="ts">
import { ref } from 'vue'
import { useTelemetryStore } from './stores/telemetryStore'
import ConceptSpeedometer from './components/ConceptSpeedometer.vue'
import ConceptAnalytics from './components/ConceptAnalytics.vue'
import ConceptChassis from './components/ConceptChassis.vue'
import ConceptStatistics from './components/ConceptStatistics.vue'
import { Wifi, WifiOff, Radio, Power } from 'lucide-vue-next'

const { status, isLiveStreamActive } = useTelemetryStore()

const isShuttingDown = ref(false)

const shutdownServer = async () => {
  if (!confirm('Are you sure you want to stop the Forza Telemetry Dash & Relay?')) return
  isShuttingDown.value = true
  try {
    await fetch('/api/shutdown', { method: 'POST' })
  } catch {
    // Ignored as server closes
  }
}
</script>

<template>
  <main class="min-h-screen bg-[#070914] text-slate-100 p-4 md:p-6 font-sans antialiased selection:bg-cyan-500/30 selection:text-cyan-200 flex flex-col justify-between">
    <div class="max-w-[1720px] mx-auto w-full space-y-6">
      <!-- Cyber Cockpit Header Bar -->
      <header class="glass-panel rounded-2xl px-6 py-4 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-center space-x-4">
          <div class="w-3.5 h-10 bg-gradient-to-b from-cyan-400 via-indigo-500 to-pink-500 rounded-sm shadow-[0_0_16px_rgba(56,189,248,0.7)]" />
          <div>
            <div class="flex items-center space-x-2">
              <h1 class="text-xl md:text-2xl font-black font-mono tracking-tight text-white uppercase">
                FORZA TELEMETRY
              </h1>
              <span class="text-[10px] px-2 py-0.5 rounded bg-cyan-950/80 text-cyan-300 border border-cyan-700/80 font-mono font-bold tracking-widest uppercase">
                CYBER GT
              </span>
            </div>
            <p class="text-[11px] text-slate-400 font-mono mt-0.5">
              CAN-BUS REALTIME TELEMETRY · 60 HZ LITTLE-ENDIAN UDP STREAM
            </p>
          </div>
        </div>

        <!-- Connection and Stream Indicators -->
        <div class="flex items-center space-x-3">
          <div
            class="flex items-center space-x-2 px-3.5 py-1.5 rounded-xl border text-xs font-mono font-bold shadow-md transition-all duration-150"
            :class="[
              isLiveStreamActive
                ? 'bg-emerald-950/50 border-emerald-500 text-emerald-300 shadow-[0_0_14px_rgba(16,185,129,0.35)] animate-pulse'
                : 'bg-slate-900/80 border-slate-800 text-slate-400'
            ]"
          >
            <Radio class="w-4 h-4" :class="isLiveStreamActive ? 'text-emerald-400' : 'text-slate-500'" />
            <span>{{ isLiveStreamActive ? 'DATA STREAM: RECEIVING' : 'DATA STREAM: WAITING FOR FORZA' }}</span>
          </div>

          <div
            class="hidden md:flex items-center space-x-2 px-3.5 py-1.5 rounded-xl border text-xs font-mono font-bold shadow-md transition-all duration-150"
            :class="[
              status === 'OPEN'
                ? 'bg-slate-900/80 border-slate-700/80 text-cyan-300'
                : 'bg-rose-950/40 border-rose-600/70 text-rose-400'
            ]"
          >
            <component :is="status === 'OPEN' ? Wifi : WifiOff" class="w-4 h-4" />
            <span>{{ status === 'OPEN' ? 'WS BRIDGE: ONLINE' : 'WS BRIDGE: OFFLINE' }}</span>
          </div>

          <!-- Quick Exit / Stop Daemon Button -->
          <button
            @click="shutdownServer"
            :disabled="isShuttingDown"
            class="flex items-center space-x-1.5 px-3 py-1.5 rounded-xl border border-rose-900/60 bg-rose-950/40 hover:bg-rose-900/60 active:scale-95 text-rose-300 text-xs font-mono font-bold transition-all duration-150 cursor-pointer shadow-md"
            title="Stop Relay & Exit Dashboard"
          >
            <Power class="w-3.5 h-3.5 text-rose-400" />
            <span>{{ isShuttingDown ? 'STOPPING...' : 'EXIT' }}</span>
          </button>
        </div>
      </header>

      <!-- Exact 4-Card Cyber Layout Matching forzadashconcept.png -->
      <section class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-4 gap-6 items-stretch">
        <!-- 1. Speedometer Pod: Dial, Gear, Rainbow Boost, kW/Nm -->
        <ConceptSpeedometer />

        <!-- 2. Dashboard Analytics: Spectrum Equalizer, Pedal Pucks, Run of Km, Tire Matrix -->
        <ConceptAnalytics />

        <!-- 3. Chassis Twin: Perspective Grid, Sculpted Supercar Silhouette, Corner Tire Slip & Dynamics -->
        <ConceptChassis />

        <!-- 4. Statistics & Dynamics: Live Wave Harmonics, Integrated G-Force Radar & Drift Horizon -->
        <ConceptStatistics />
      </section>

      <!-- Footer System Status -->
      <footer class="glass-panel rounded-xl px-5 py-2.5 flex flex-col sm:flex-row items-center justify-between text-[11px] font-mono text-slate-500 gap-2">
        <div class="flex items-center space-x-2">
          <span class="w-2 h-2 rounded-full bg-cyan-400 animate-pulse" />
          <span>CYBER DASHBOARD TELEMETRY BRIDGE · ACTIVE UDP PORT 2048</span>
        </div>
        <div class="flex items-center space-x-5 font-bold text-slate-400">
          <span>SAMPLING: 60 FPS</span>
          <span>LATENCY: ZERO-COPY</span>
          <span>UNITS: METRIC (KM/H, °C, KW, N·M, BAR)</span>
        </div>
      </footer>
    </div>
  </main>
</template>
