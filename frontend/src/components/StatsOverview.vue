<script setup lang="ts">
import { Activity, Bell, ShieldAlert, Zap } from 'lucide-vue-next';

defineProps<{
  stats: any;
}>();
</script>

<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
    <div class="glass p-8 ios-card flex items-center gap-6 group">
      <div class="w-16 h-16 rounded-[1.5rem] bg-ios-accent/15 text-ios-accent flex items-center justify-center shadow-inner group-hover:scale-110 transition-transform duration-700">
        <Activity :size="32" />
      </div>
      <div>
        <div class="text-[11px] font-black text-ios-secondary uppercase tracking-[0.2em] mb-2">Throughput</div>
        <div class="text-4xl font-black text-white tracking-tighter tabular-nums">{{ stats.total_packets.toLocaleString() }}</div>
      </div>
    </div>

    <div class="glass p-8 ios-card flex items-center gap-6 group">
      <div class="w-16 h-16 rounded-[1.5rem] bg-ios-danger/15 text-ios-danger flex items-center justify-center shadow-inner group-hover:scale-110 transition-transform duration-700">
        <Bell :size="32" />
      </div>
      <div>
        <div class="text-[11px] font-black text-ios-secondary uppercase tracking-[0.2em] mb-2">Security</div>
        <div class="text-4xl font-black text-white tracking-tighter tabular-nums">{{ stats.total_alerts.toLocaleString() }}</div>
      </div>
    </div>

    <div class="glass p-8 ios-card flex items-center gap-6 group">
      <div class="w-16 h-16 rounded-[1.5rem] bg-ios-success/15 text-ios-success flex items-center justify-center shadow-inner group-hover:scale-110 transition-transform duration-700">
        <ShieldCheck v-if="stats.total_alerts === 0" :size="32" />
        <ShieldAlert v-else :size="32" class="text-ios-warning" />
      </div>
      <div>
        <div class="text-[11px] font-black text-ios-secondary uppercase tracking-[0.2em] mb-2">Status</div>
        <div :class="['text-3xl font-black tracking-tighter', stats.total_alerts > 0 ? 'text-ios-warning' : 'text-ios-success']">
          {{ stats.total_alerts > 0 ? 'Elevated' : 'Shielded' }}
        </div>
      </div>
    </div>

    <div class="glass p-8 ios-card flex items-center gap-6 group">
      <div class="w-16 h-16 rounded-[1.5rem] bg-purple-500/15 text-purple-400 flex items-center justify-center shadow-inner group-hover:scale-110 transition-transform duration-700">
        <Zap :size="32" />
      </div>
      <div>
        <div class="text-[11px] font-black text-ios-secondary uppercase tracking-[0.2em] mb-2">Protocols</div>
        <div class="text-4xl font-black text-white tracking-tighter tabular-nums">{{ Object.keys(stats.protocol_dist).length }}</div>
      </div>
    </div>
  </div>
</template>
