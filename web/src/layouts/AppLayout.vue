<script setup lang="ts">
import AppTopbar from '@/components/AppTopbar.vue'
</script>

<template>
  <div class="app-layout">
    <AppTopbar />
    <main class="app-layout__main">
      <router-view />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: var(--color-bg);
}

.app-layout__main {
  min-height: calc(100vh - 50px);
  background: var(--color-bg);
}
</style>
