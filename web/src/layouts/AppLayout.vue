<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import AppContextPanel from '@/components/AppContextPanel.vue'
import AppModuleSidebar from '@/components/AppModuleSidebar.vue'
import AppTopbar from '@/components/AppTopbar.vue'
import CreateProjectDialog from '@/components/CreateProjectDialog.vue'
import CreateWorkspaceDialog from '@/components/CreateWorkspaceDialog.vue'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'

const route = useRoute()
const { bootstrap } = useWorkspaceContext()
const showOverlay = computed(() => route.meta.overlay === true)

onMounted(() => {
  bootstrap()
})
</script>

<template>
  <div class="app-layout">
    <AppTopbar />
    <div class="app-layout__body">
      <div class="app-layout__nav" :class="{ 'app-layout__nav--dimmed': showOverlay }">
        <AppModuleSidebar />
        <AppContextPanel />
      </div>

      <main class="app-layout__main">
        <router-view />
      </main>

      <router-view v-slot="{ Component }" name="overlay">
        <Transition name="app-layout__overlay">
          <component :is="Component" v-if="Component" class="app-layout__overlay" />
        </Transition>
      </router-view>
    </div>

    <CreateWorkspaceDialog />
    <CreateProjectDialog />
  </div>
</template>

<style scoped>
.app-layout {
  min-height: 100vh;
  background: var(--color-bg);
  display: flex;
  flex-direction: column;
}

.app-layout__body {
  position: relative;
  flex: 1;
  display: flex;
  min-height: 0;
  height: calc(100vh - 56px);
  overflow: hidden;
}

.app-layout__nav {
  display: flex;
  flex-shrink: 0;
  transition: opacity 0.2s ease;
}

.app-layout__nav--dimmed {
  opacity: 0.35;
  pointer-events: none;
}

.app-layout__main {
  flex: 1;
  min-width: 0;
  overflow: hidden;
  background: var(--color-bg);
}

.app-layout__overlay {
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  width: min(760px, 62vw);
  z-index: 120;
  background: var(--color-bg);
  border-right: 1px solid var(--color-border);
  box-shadow: 8px 0 32px rgba(0, 0, 0, 0.18);
  overflow: hidden;
}
</style>
