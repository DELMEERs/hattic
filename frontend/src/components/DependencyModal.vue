<script setup lang="ts">
import { AlertCircle, Download, Terminal, ShieldAlert } from 'lucide-vue-next';

defineProps<{
  status: any;
}>();

const openUrl = (url: string) => {
  window.open(url, '_blank');
};

const reloadPage = () => {
  window.location.reload();
};
</script>

<template>
  <transition
    enter-active-class="transition duration-500 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition duration-300 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div v-if="status && status.status === 'ERROR'" class="fixed inset-0 z-[200] flex items-center justify-center p-6 bg-black/60 backdrop-blur-xl">
      <div class="w-full max-w-md glass rounded-[2.5rem] ios-shadow p-8 text-center space-y-8 border border-white/20">
        <div class="flex justify-center">
          <div class="w-24 h-24 rounded-[2rem] bg-ios-danger/20 flex items-center justify-center shadow-2xl animate-bounce">
            <ShieldAlert class="text-ios-danger" :size="48" />
          </div>
        </div>

        <div class="space-y-3">
          <h2 class="text-3xl font-black tracking-tighter text-white">Dependency Missing</h2>
          <p class="text-ios-secondary text-sm font-medium leading-relaxed">
            Hattic requires specific system-level components to monitor network traffic. 
            Please follow the instructions below for <span class="text-white font-bold">{{ status.platform }}</span>.
          </p>
        </div>

        <div class="p-6 rounded-3xl bg-white/5 border border-white/10 text-left space-y-4">
          <div class="flex items-center gap-3">
            <div class="w-8 h-8 rounded-lg bg-ios-accent/20 flex items-center justify-center">
              <AlertCircle :size="18" class="text-ios-accent" />
            </div>
            <span class="text-xs font-black uppercase tracking-widest text-white/80">Requirement: {{ status.missingDep }}</span>
          </div>
          
          <div class="space-y-2">
            <div class="text-[10px] font-black text-ios-secondary uppercase tracking-widest px-1">Instructions</div>
            <div v-if="status.platform === 'linux'" class="bg-black/40 p-4 rounded-xl font-mono text-[11px] text-ios-success break-all leading-relaxed border border-ios-success/20 flex items-start gap-3">
              <Terminal :size="14" class="shrink-0 mt-0.5" />
              <code>{{ status.instructions }}</code>
            </div>
            <p v-else class="text-xs text-white/70 font-medium px-1 leading-relaxed">
              {{ status.instructions }}
            </p>
          </div>
        </div>

        <div class="flex flex-col gap-3">
          <button 
            v-if="status.platform === 'windows'"
            @click="openUrl('https://npcap.com/')"
            class="w-full py-4 rounded-2xl bg-ios-accent text-white font-black text-sm uppercase tracking-[0.2em] shadow-xl shadow-ios-accent/20 hover:opacity-90 active:scale-95 transition-all flex items-center justify-center gap-3"
          >
            <Download :size="18" />
            Download Npcap
          </button>
          
          <button 
            @click="reloadPage"
            class="w-full py-4 rounded-2xl bg-white/5 text-white font-black text-sm uppercase tracking-[0.2em] hover:bg-white/10 transition-all border border-white/10"
          >
            Check Again
          </button>
        </div>
      </div>
    </div>
  </transition>
</template>