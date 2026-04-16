<script setup lang="ts">
import { ref, onMounted } from "vue";

const status = ref("this used to be just some janky school project");
const dots = ref("");

onMounted(() => {
    setInterval(() => {
        dots.value = dots.value.length >= 3 ? "" : dots.value + ".";
    }, 500);
});
</script>

<template>
    <div class="gradient-bg"></div>
    <div class="creative-placeholder">
        <div class="logo"></div>
        <h1 class="title">hattic</h1>
        <div class="status-text">{{ status }}{{ dots }}</div>
    </div>
</template>

<style scoped></style>
