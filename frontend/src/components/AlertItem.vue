<script setup lang="ts">
import { AlertCircle, AlertTriangle, Info } from 'lucide-vue-next';

defineProps<{
  alert: any;
}>();

const getLevelStyles = (level: string) => {
  switch (level) {
    case 'Critical': return 'from-red-500 to-pink-500 text-white shadow-red-500/20';
    case 'Warning': return 'from-orange-400 to-yellow-500 text-white shadow-orange-500/20';
    default: return 'from-blue-400 to-ios-accent text-white shadow-blue-500/20';
  }
};

const getIcon = (level: string) => {
  switch (level) {
    case 'Critical': return AlertCircle;
    case 'Warning': return AlertTriangle;
    default: return Info;
  }
};
</script>

<template>
  <div :class="['px-6 py-5 border-b border-white/5 hover:bg-white/5 transition-all duration-300 group', alert.level === 'Critical' ? 'bg-red-500/5 glow-critical my-2 rounded-3xl border-transparent mx-2' : '']">
    <div class="flex gap-5">
      <div :class="['w-14 h-14 rounded-[1.25rem] flex items-center justify-center bg-gradient-to-br shadow-xl shrink-0 transition-transform duration-500 group-hover:scale-110', getLevelStyles(alert.level)]">
        <component :is="getIcon(alert.level)" :size="28" />
      </div>
      
      <div class="flex-1">
        <div class="flex items-center justify-between mb-1.5">
          <span :class="['text-sm font-bold tracking-tight', alert.level === 'Critical' ? 'text-red-400' : 'text-white']">{{ alert.type }}</span>
          <span class="text-[10px] text-ios-secondary uppercase font-bold tracking-widest opacity-60">
            {{ new Date(alert.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' }) }}
          </span>
        </div>
        <p class="text-sm text-ios-secondary leading-relaxed font-medium group-hover:text-white/90 transition-colors">
          {{ alert.message }}
        </p>
        <div v-if="alert.src_ip" class="mt-3 inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/5 border border-white/10 text-[10px] text-ios-accent font-bold uppercase tracking-wider">
          <div class="w-1 h-1 rounded-full bg-ios-accent animate-pulse"></div>
          Source: {{ alert.src_ip }}
        </div>
      </div>
    </div>
  </div>
</template>
