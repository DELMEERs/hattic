<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import Sidebar from './components/Sidebar.vue';
import StatsOverview from './components/StatsOverview.vue';
import GlassCard from './components/GlassCard.vue';
import TrafficItem from './components/TrafficItem.vue';
import AlertItem from './components/AlertItem.vue';
import ConfigPanel from './components/ConfigPanel.vue';
import { EventsOn, EventsOff } from '../wailsjs/runtime';
import { GetStats, GetIsSniffing, HealthCheck } from '../wailsjs/go/main/App';
import { Trash2, Activity, Bell, ChevronDown, Zap, ShieldCheck, ShieldAlert } from 'lucide-vue-next';

const activeTab = ref('dashboard');
const stats = ref({ total_packets: 0, total_alerts: 0, protocol_dist: {} });
const recentPackets = ref<any[]>([]);
const alerts = ref<any[]>([]);
const isSniffing = ref(false);
const showLogPopover = ref(false);
const health = ref({ is_root: false, pcap_version: '', can_sniff: false, error: '' });
const lastError = ref('');

const updateStats = async () => {
  try {
    stats.value = await GetStats();
    isSniffing.value = await GetIsSniffing();
    health.value = await HealthCheck();
  } catch (err) {
    console.error(err);
  }
};

onMounted(() => {
  updateStats();
  const statsInterval = setInterval(updateStats, 2000);

  EventsOn('packet', (packet: any) => {
    recentPackets.value.unshift(packet);
    if (recentPackets.value.length > 50) {
      recentPackets.value.pop();
    }
  });

  EventsOn('alert', (alert: any) => {
    alerts.value.unshift(alert);
    if (alerts.value.length > 100) {
      alerts.value.pop();
    }
  });

  EventsOn('sniffer_status', (status: boolean) => {
    isSniffing.value = status;
    if (status) lastError.value = '';
  });

  EventsOn('backend-error', (err: string) => {
    lastError.value = err;
  });

  onUnmounted(() => {
    clearInterval(statsInterval);
    EventsOff('packet');
    EventsOff('alert');
    EventsOff('sniffer_status');
    EventsOff('backend-error');
  });
});

const clearTraffic = () => {
  recentPackets.value = [];
};

const clearAlerts = () => {
  alerts.value = [];
};
</script>

<template>
  <!-- Animated Mesh Background -->
  <div class="mesh-gradient">
    <div class="mesh-blob w-[600px] h-[600px] bg-purple-600 top-[-10%] left-[-10%]"></div>
    <div class="mesh-blob w-[500px] h-[500px] bg-blue-600 bottom-[-10%] right-[-10%] delay-700"></div>
    <div class="mesh-blob w-[400px] h-[400px] bg-ios-accent top-[20%] right-[10%] delay-1000"></div>
  </div>

  <div class="flex h-screen w-full text-white overflow-hidden relative z-10 font-sans">
    <Sidebar v-model:activeTab="activeTab" />
    
    <main class="flex-1 flex flex-col min-w-0 h-full">
      <!-- Header Area -->
      <header class="h-24 flex items-center justify-between px-10 shrink-0 relative z-50">
        <div>
          <h2 class="text-4xl font-black tracking-tighter text-white/95">{{ activeTab }}</h2>
          <div class="flex items-center gap-2 mt-1">
            <span class="text-ios-secondary text-[10px] font-black uppercase tracking-widest">Network Sentinel v1.0</span>
            <span class="w-1 h-1 rounded-full bg-white/20"></span>
            <span class="text-ios-secondary text-[10px] font-black uppercase tracking-widest">{{ health.pcap_version }}</span>
          </div>
        </div>
        
        <div class="flex items-center gap-6">
          <!-- Connection Status Badge -->
          <div :class="['px-5 py-2.5 glass rounded-2xl flex items-center gap-4 transition-all duration-500', 
                        isSniffing ? 'border-ios-success/40 ring-4 ring-ios-success/10' : (lastError || !health.can_sniff ? 'border-ios-danger/40 ring-4 ring-ios-danger/10' : 'border-white/5')]">
            <div :class="['w-3 h-3 rounded-full shadow-[0_0_10px_rgba(52,199,89,0.5)]', isSniffing ? 'bg-ios-success animate-pulse' : (lastError || !health.can_sniff ? 'bg-ios-danger shadow-[0_0_10px_rgba(255,59,48,0.5)]' : 'bg-ios-warning')]"></div>
            <div class="flex flex-col">
              <span class="text-[10px] font-black uppercase tracking-widest leading-none">
                {{ isSniffing ? 'ACTIVE' : (lastError || !health.can_sniff ? 'ERROR' : 'IDLE') }}
              </span>
              <span v-if="lastError || health.error" class="text-[9px] text-ios-danger font-bold mt-1.5 max-w-[180px] truncate">
                {{ lastError || health.error }}
              </span>
            </div>
          </div>

          <div 
            @click="showLogPopover = !showLogPopover"
            class="p-3.5 glass rounded-2xl cursor-pointer hover:bg-white/10 transition-all active:scale-95 group"
          >
            <Activity :size="22" :class="[isSniffing ? 'text-ios-success animate-pulse' : 'text-white/30 group-hover:text-white/60', 'transition-colors']" />
          </div>

          <!-- Log Popover -->
          <transition 
            enter-active-class="transition duration-400 ease-out"
            enter-from-class="translate-y-4 opacity-0 scale-95"
            enter-to-class="translate-y-0 opacity-100 scale-100"
            leave-active-class="transition duration-300 ease-in"
            leave-from-class="translate-y-0 opacity-100 scale-100"
            leave-to-class="translate-y-4 opacity-0 scale-95"
          >
            <div v-if="showLogPopover" class="absolute top-[90%] right-10 w-96 glass-dark rounded-[2.5rem] ios-shadow p-6 z-[100] border border-white/20">
              <div class="flex items-center justify-between mb-6 px-1">
                <span class="text-[11px] font-black uppercase tracking-[0.2em] text-ios-secondary">Real-Time Capture</span>
                <div class="flex items-center gap-2">
                   <div class="w-2 h-2 rounded-full bg-ios-success animate-pulse"></div>
                   <Activity :size="16" class="text-ios-success" />
                </div>
              </div>
              <div class="space-y-3 max-h-[400px] overflow-y-auto pr-2 custom-scrollbar">
                <div v-if="recentPackets.length === 0" class="py-16 text-center text-ios-secondary text-sm opacity-40 flex flex-col items-center gap-4">
                  <div class="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center">
                    <Activity :size="32" />
                  </div>
                  <p class="font-bold tracking-tight">Initializing sensor array...</p>
                </div>
                <div v-for="(p, i) in recentPackets.slice(0, 8)" :key="i" class="flex items-center gap-4 p-3.5 rounded-2xl bg-white/5 border border-white/5 hover:bg-white/10 transition-all group">
                  <div :class="['w-10 h-10 rounded-xl flex items-center justify-center shrink-0 shadow-inner', p.protocol === 'TCP' ? 'bg-ios-accent/20 text-ios-accent' : 'bg-purple-500/20 text-purple-400']">
                    <Zap :size="16" />
                  </div>
                  <div class="flex-1 min-w-0">
                    <div class="text-xs font-black text-white flex items-center justify-between">
                      <span class="tracking-tight">{{ p.protocol }}</span>
                      <span class="text-[10px] text-ios-secondary tabular-nums opacity-60">{{ p.length }}B</span>
                    </div>
                    <div class="text-[10px] text-ios-secondary truncate mt-1 font-medium opacity-80">
                      {{ p.src_ip }} <span class="text-ios-accent">→</span> {{ p.dst_ip }}
                    </div>
                  </div>
                </div>
              </div>
              <div class="mt-6 pt-5 border-t border-white/10 flex justify-center">
                <button @click="activeTab = 'traffic'; showLogPopover = false" class="text-[11px] font-black text-ios-accent uppercase tracking-[0.2em] hover:text-white transition-all">
                  Access Secure Vault
                </button>
              </div>
            </div>
          </transition>
        </div>
      </header>

      <!-- Content Area -->
      <div class="flex-1 overflow-auto p-10 space-y-10 scroll-smooth custom-scrollbar">
        
        <template v-if="activeTab === 'dashboard'">
          <StatsOverview :stats="stats" />
          
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <GlassCard title="Recent Activity" class="h-[500px]">
              <template #header-action>
                <button @click="clearTraffic" class="p-2 hover:bg-white/10 rounded-full transition-colors text-ios-secondary">
                  <Trash2 :size="18" />
                </button>
              </template>
              <div v-if="recentPackets.length === 0" class="h-full flex flex-col items-center justify-center text-ios-secondary opacity-50 space-y-4">
                <Activity :size="48" />
                <p>Waiting for network traffic...</p>
              </div>
              <TrafficItem v-for="(p, i) in recentPackets.slice(0, 10)" :key="i" :packet="p" />
            </GlassCard>

            <GlassCard title="Security Alerts" class="h-[500px]">
              <template #header-action>
                <button @click="clearAlerts" class="p-2 hover:bg-white/10 rounded-full transition-colors text-ios-secondary">
                  <Trash2 :size="18" />
                </button>
              </template>
              <div v-if="alerts.length === 0" class="h-full flex flex-col items-center justify-center text-ios-secondary opacity-50 space-y-4">
                <Bell :size="48" />
                <p>No threats detected</p>
              </div>
              <AlertItem v-for="(a, i) in alerts.slice(0, 10)" :key="i" :alert="a" />
            </GlassCard>
          </div>
        </template>

        <template v-if="activeTab === 'traffic'">
          <GlassCard title="Live Packet Stream" class="min-h-[600px]">
            <template #header-action>
              <span class="text-xs font-medium text-ios-secondary">{{ recentPackets.length }} packets buffered</span>
            </template>
            <div v-if="recentPackets.length === 0" class="h-[400px] flex flex-col items-center justify-center text-ios-secondary opacity-50 space-y-4">
              <Activity :size="48" />
              <p>Waiting for network traffic...</p>
            </div>
            <TrafficItem v-for="(p, i) in recentPackets" :key="i" :packet="p" />
          </GlassCard>
        </template>

        <template v-if="activeTab === 'alerts'">
          <GlassCard title="Alert Center" class="min-h-[600px]">
            <template #header-action>
              <span class="text-xs font-medium text-ios-secondary">{{ alerts.length }} alerts recorded</span>
            </template>
            <div v-if="alerts.length === 0" class="h-[400px] flex flex-col items-center justify-center text-ios-secondary opacity-50 space-y-4">
              <Bell :size="48" />
              <p>Your network is safe</p>
            </div>
            <div class="space-y-0">
              <AlertItem v-for="(a, i) in alerts" :key="i" :alert="a" />
            </div>
          </GlassCard>
        </template>

        <template v-if="activeTab === 'settings'">
          <ConfigPanel />
        </template>

      </div>
    </main>
  </div>
</template>

<style>
/* Custom scrollbar for webkit */
::-webkit-scrollbar {
  width: 8px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 10px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.2);
}
</style>
