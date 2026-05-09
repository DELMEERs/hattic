<script setup lang="ts">
import { ArrowRight, Globe, Shield, Zap } from 'lucide-vue-next';

defineProps<{
  packet: any;
}>();

const getProtocolColor = (proto: string) => {
  switch (proto) {
    case 'TCP': return 'bg-blue-500/20 text-blue-400';
    case 'UDP': return 'bg-purple-500/20 text-purple-400';
    case 'HTTP': return 'bg-green-500/20 text-green-400';
    case 'HTTPS': return 'bg-emerald-500/20 text-emerald-400';
    case 'DNS': return 'bg-yellow-500/20 text-yellow-400';
    case 'MDNS': return 'bg-orange-500/20 text-orange-400';
    case 'SSH': return 'bg-slate-500/20 text-slate-400';
    default: return 'bg-gray-500/20 text-gray-400';
  }
};
</script>

<template>
  <div class="px-6 py-4 flex items-center gap-5 border-b border-white/5 hover:bg-white/10 transition-all duration-300 cursor-default group">
    <div :class="['w-12 h-12 rounded-2xl flex items-center justify-center shadow-lg transition-transform duration-500 group-hover:scale-110', getProtocolColor(packet.protocol)]">
      <Zap v-if="packet.protocol === 'TCP' || packet.protocol === 'UDP'" :size="20" />
      <Globe v-else-if="packet.protocol === 'HTTP' || packet.protocol === 'HTTPS'" :size="20" />
      <Shield v-else :size="20" />
    </div>
    
    <div class="flex-1 min-w-0">
      <div class="flex items-center gap-3">
        <span class="text-[13px] font-bold text-white tracking-tight tabular-nums">{{ packet.src_ip }}</span>
        <div class="flex items-center gap-1 opacity-40">
           <div class="w-1 h-1 rounded-full bg-white"></div>
           <div class="w-1.5 h-1.5 rounded-full bg-white"></div>
           <ArrowRight :size="12" class="text-white" />
        </div>
        <span class="text-[13px] font-bold text-white tracking-tight tabular-nums">{{ packet.dst_ip }}</span>
      </div>
      <div class="flex items-center gap-2.5 mt-1.5">
        <span class="text-[10px] font-black text-ios-accent uppercase tracking-widest bg-ios-accent/10 px-2 py-0.5 rounded-md">{{ packet.protocol }}</span>
        <span class="text-[10px] text-ios-secondary font-bold">•</span>
        <span class="text-[10px] text-ios-secondary font-bold opacity-70">{{ packet.length }} BYTES</span>
        <span v-if="packet.hostname" class="text-[10px] text-ios-success font-black truncate uppercase tracking-wider">• {{ packet.hostname }}</span>
      </div>
    </div>

    <div class="text-right hidden sm:block">
      <div class="text-[11px] font-black text-white/80 tabular-nums tracking-widest">{{ packet.src_port }} <span class="opacity-30">:</span> {{ packet.dst_port }}</div>
      <div class="text-[9px] font-bold text-ios-secondary mt-1 opacity-50 uppercase tracking-tighter">
        {{ new Date(packet.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }) }}
      </div>
    </div>
  </div>
</template>
