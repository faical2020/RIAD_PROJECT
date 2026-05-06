<template>
    <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import syncProvider from './sync'

onMounted(async () => {
    console.log('[Sync] Mounting App.vue...');
    await syncProvider.init();
})
</script>


<style>
body {
    margin: 0;
    padding: 0;
    font-family: system-ui, -apple-system, sans-serif;
}
</style>
