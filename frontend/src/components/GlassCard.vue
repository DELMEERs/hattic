<script setup lang="ts">
defineProps<{
  title?: string;
  class?: string;
}>();
</script>

<template>
  <div :class="['glass ios-card overflow-hidden flex flex-col transition-all duration-300', $props.class]">
    <div v-if="title" class="px-6 py-4 border-b border-white/10 flex items-center justify-between">
      <h3 class="text-lg font-bold tracking-tight text-white/90">{{ title }}</h3>
      <slot name="header-action"></slot>
    </div>
    <div class="flex-1 overflow-auto">
      <slot></slot>
    </div>
  </div>
</template>
