<script setup lang="ts">
import { Activity, Bell, Settings, LayoutDashboard } from 'lucide-vue-next';

defineProps<{
  activeTab: string;
}>();

defineEmits(['update:activeTab']);
</script>

<template>
  <div class="w-24 lg:w-80 h-full glass-dark border-r border-white/5 flex flex-col py-12 px-6 z-20">
    <div class="flex items-center gap-5 px-4 mb-16 group">
      <div class="w-14 h-14 rounded-[1.5rem] bg-gradient-to-tr from-ios-accent to-blue-500 flex items-center justify-center shadow-xl shadow-ios-accent/20 group-hover:rotate-[360deg] transition-transform duration-1000">
        <Activity class="text-white" :size="32" />
      </div>
      <div class="hidden lg:block">
        <h1 class="text-3xl font-black tracking-tighter text-white">HATTIC</h1>
        <div class="text-[10px] font-black text-ios-accent uppercase tracking-[0.3em] opacity-80 mt-1">Network Guardian</div>
      </div>
    </div>

    <nav class="flex-1 space-y-4">
      <button 
        @click="$emit('update:activeTab', 'dashboard')"
        :class="['w-full flex items-center gap-5 px-6 py-5 rounded-[1.5rem] transition-all duration-500 group relative overflow-hidden', 
                 activeTab === 'dashboard' ? 'glass text-white shadow-xl' : 'text-ios-secondary hover:text-white hover:bg-white/5']"
      >
        <div v-if="activeTab === 'dashboard'" class="absolute inset-0 bg-ios-accent/10 -z-10"></div>
        <LayoutDashboard :size="28" :class="['transition-all duration-500', activeTab === 'dashboard' ? 'text-ios-accent scale-110' : 'group-hover:scale-110']" />
        <span class="font-black tracking-tight hidden lg:block uppercase text-xs tracking-[0.1em]">Dashboard</span>
      </button>

      <button 
        @click="$emit('update:activeTab', 'traffic')"
        :class="['w-full flex items-center gap-5 px-6 py-5 rounded-[1.5rem] transition-all duration-500 group relative overflow-hidden', 
                 activeTab === 'traffic' ? 'glass text-white shadow-xl' : 'text-ios-secondary hover:text-white hover:bg-white/5']"
      >
        <div v-if="activeTab === 'traffic'" class="absolute inset-0 bg-ios-accent/10 -z-10"></div>
        <Activity :size="28" :class="['transition-all duration-500', activeTab === 'traffic' ? 'text-ios-accent scale-110' : 'group-hover:scale-110']" />
        <span class="font-black tracking-tight hidden lg:block uppercase text-xs tracking-[0.1em]">Live Stream</span>
      </button>

      <button 
        @click="$emit('update:activeTab', 'alerts')"
        :class="['w-full flex items-center gap-5 px-6 py-5 rounded-[1.5rem] transition-all duration-500 group relative overflow-hidden', 
                 activeTab === 'alerts' ? 'glass text-white shadow-xl' : 'text-ios-secondary hover:text-white hover:bg-white/5']"
      >
        <div v-if="activeTab === 'alerts'" class="absolute inset-0 bg-ios-accent/10 -z-10"></div>
        <Bell :size="28" :class="['transition-all duration-500', activeTab === 'alerts' ? 'text-ios-accent scale-110' : 'group-hover:scale-110']" />
        <span class="font-black tracking-tight hidden lg:block uppercase text-xs tracking-[0.1em]">Vault</span>
      </button>
    </nav>

    <div class="mt-auto pt-10 border-t border-white/5">
      <button 
        @click="$emit('update:activeTab', 'settings')"
        :class="['w-full flex items-center gap-5 px-6 py-5 rounded-[1.5rem] transition-all duration-500 group relative overflow-hidden', 
                 activeTab === 'settings' ? 'glass text-white shadow-xl' : 'text-ios-secondary hover:text-white hover:bg-white/5']"
      >
        <div v-if="activeTab === 'settings'" class="absolute inset-0 bg-ios-accent/10 -z-10"></div>
        <Settings :size="28" :class="['transition-all duration-500', activeTab === 'settings' ? 'text-ios-accent scale-110' : 'group-hover:scale-110']" />
        <span class="font-black tracking-tight hidden lg:block uppercase text-xs tracking-[0.1em]">Settings</span>
      </button>
    </div>
  </div>
</template>
