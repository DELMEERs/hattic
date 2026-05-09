<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { GetInterfaces, GetConfig, SaveConfig, StartSniffer, StopSniffer, GetIsSniffing, TriggerTestAlert } from '../../wailsjs/go/main/App';
import GlassCard from './GlassCard.vue';
import { Play, Square, Settings, Wifi, ShieldAlert } from 'lucide-vue-next';

const interfaces = ref<any[]>([]);
const config = ref<any>({ interface_name: '', promiscuous: true, snap_len: 65536 });
const isSniffing = ref(false);
const loading = ref(true);

onMounted(async () => {
  try {
    interfaces.value = await GetInterfaces();
    config.value = await GetConfig();
    isSniffing.value = await GetIsSniffing();
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
});

const handleSave = async () => {
  try {
    await SaveConfig(config.value);
  } catch (err) {
    console.error(err);
  }
};

const toggleSniffer = async () => {
  try {
    if (isSniffing.value) {
      await StopSniffer();
    } else {
      await StartSniffer();
    }
    isSniffing.value = !isSniffing.value;
  } catch (err) {
    console.error(err);
  }
};

const triggerTestAnomaly = async () => {
  try {
    await TriggerTestAlert();
  } catch (err) {
    console.error(err);
  }
};
</script>

<template>
  <div class="space-y-10">
    <GlassCard title="Sensor Configuration">
      <div class="p-10 space-y-8">
        <div class="space-y-4">
          <label class="text-[11px] font-black text-ios-secondary uppercase tracking-[0.2em] px-2">Primary Inbound Interface</label>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
            <button 
              v-for="iface in interfaces" 
              :key="iface.name"
              @click="config.interface_name = iface.name"
              :class="['text-left p-6 rounded-[1.75rem] border transition-all duration-500 flex items-start gap-4 group', 
                       config.interface_name === iface.name ? 'bg-ios-accent/10 border-ios-accent ring-2 ring-ios-accent/20' : 'bg-white/5 border-white/5 hover:border-white/20']"
            >
              <div :class="['w-10 h-10 rounded-xl flex items-center justify-center shrink-0 transition-transform duration-500 group-hover:scale-110', 
                           config.interface_name === iface.name ? 'bg-ios-accent text-white shadow-lg' : 'bg-white/10 text-ios-secondary']">
                <Wifi :size="20" />
              </div>
              <div class="min-w-0">
                <div class="text-[14px] font-black text-white truncate tracking-tight">{{ iface.name }}</div>
                <div class="text-[10px] text-ios-secondary truncate mt-1 font-bold">{{ iface.description || 'GENERIC ADAPTER' }}</div>
                <div v-if="iface.ip_addresses" class="text-[10px] text-ios-accent mt-2 font-black tabular-nums">{{ iface.ip_addresses[0] }}</div>
              </div>
            </button>
          </div>
        </div>

        <div class="flex items-center gap-10 pt-4">
          <label class="flex items-center gap-4 cursor-pointer group">
            <div class="relative">
              <input type="checkbox" v-model="config.promiscuous" class="sr-only">
              <div :class="['w-14 h-8 rounded-full transition-all duration-500', config.promiscuous ? 'bg-ios-success shadow-[0_0_15px_rgba(52,199,89,0.4)]' : 'bg-white/10']"></div>
              <div :class="['absolute top-1 left-1 w-6 h-6 rounded-full bg-white transition-all duration-500 shadow-xl', config.promiscuous ? 'translate-x-6' : 'translate-x-0']"></div>
            </div>
            <div class="flex flex-col">
              <span class="text-sm font-black text-white tracking-tight">Promiscuous Mode</span>
              <span class="text-[10px] text-ios-secondary font-bold uppercase tracking-widest mt-0.5">Capture all network packets</span>
            </div>
          </label>
        </div>

        <div class="pt-6 border-t border-white/5 flex justify-end">
          <button 
            @click="handleSave"
            class="px-8 py-4 rounded-2xl bg-ios-accent text-white text-xs font-black uppercase tracking-[0.2em] shadow-xl shadow-ios-accent/20 hover:opacity-90 active:scale-95 transition-all"
          >
            Save Configuration
          </button>
        </div>
      </div>
    </GlassCard>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-10">
      <GlassCard title="Engine Control">
        <div class="p-10 flex flex-col items-center justify-center text-center space-y-8">
          <div :class="['w-28 h-28 rounded-[3rem] flex items-center justify-center shadow-2xl transition-all duration-700', 
                       isSniffing ? 'bg-ios-danger shadow-ios-danger/30 scale-110 rotate-[360deg]' : 'bg-ios-success shadow-ios-success/30 hover:scale-105']">
            <component :is="isSniffing ? Square : Play" class="text-white fill-white" :size="40" />
          </div>

          <div>
            <h4 class="text-3xl font-black text-white tracking-tighter mb-3">{{ isSniffing ? 'Engine Active' : 'Engine Ready' }}</h4>
            <p class="text-ios-secondary text-sm font-medium max-w-[280px] leading-relaxed mx-auto">
              {{ isSniffing ? 'Deep packet inspection is currently active across the selected interface.' : 'Initialize the Sentinel engine to begin real-time network analysis.' }}
            </p>
          </div>

          <button 
            @click="toggleSniffer"
            :class="['px-12 py-5 rounded-[2rem] font-black text-sm uppercase tracking-[0.2em] shadow-2xl transition-all duration-500 transform active:scale-95', 
                     isSniffing ? 'bg-white/5 text-white border border-white/10 hover:bg-white/10' : 'bg-white text-black hover:bg-white/90']"
          >
            {{ isSniffing ? 'Stop Sentinel' : 'Start Sentinel' }}
          </button>
        </div>
      </GlassCard>

      <GlassCard title="Security Lab">
        <div class="p-10 flex flex-col items-center justify-center text-center space-y-8">
          <div class="w-28 h-28 rounded-[3rem] bg-ios-warning/10 flex items-center justify-center shadow-2xl border border-ios-warning/20">
            <ShieldAlert class="text-ios-warning" :size="56" />
          </div>

          <div>
            <h4 class="text-3xl font-black text-white tracking-tighter mb-3">Anomaly Lab</h4>
            <p class="text-ios-secondary text-sm font-medium max-w-[280px] leading-relaxed mx-auto">
              Test your security configuration by injecting a simulated high-severity network anomaly.
            </p>
          </div>

          <button 
            @click="triggerTestAnomaly"
            class="px-12 py-5 rounded-[2rem] font-black text-sm uppercase tracking-[0.2em] bg-ios-danger/10 text-ios-danger border border-ios-danger/30 hover:bg-ios-danger/20 transition-all shadow-xl active:scale-95"
          >
            Inject Anomaly
          </button>
        </div>
      </GlassCard>
    </div>
  </div>
</template>

