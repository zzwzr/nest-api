<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Moon, Setting, Sunny, SwitchButton } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import LocaleSwitcher from '@/components/LocaleSwitcher.vue'
import WorkspaceProjectSelector from '@/components/WorkspaceProjectSelector.vue'
import { useAuth } from '@/composables/useAuth'
import { useLocale } from '@/composables/useLocale'
import { useTheme } from '@/composables/useTheme'

const router = useRouter()
const { user, logout } = useAuth()
const { theme, toggleTheme } = useTheme()
const { t } = useLocale()

const avatarText = computed(() => {
  const name = user.value?.name || user.value?.account || '?'
  return name.slice(0, 1).toUpperCase()
})

const displayName = computed(() => user.value?.name || user.value?.account || '')

const avatarSrc = computed(() => user.value?.avatar || '')

const themeTooltip = computed(() =>
  theme.value === 'light' ? t('topbar.themeLight') : t('topbar.themeDark'),
)

function goSettings() {
  router.push('/admin/settings')
}

function handleLogout() {
  logout()
  ElMessage.success(t('auth.logoutSuccess'))
  router.push('/login')
}

function handleCommand(command: string) {
  if (command === 'logout') {
    handleLogout()
  }
}
</script>

<template>
  <header class="app-topbar">
    <div class="app-topbar__left">
      <router-link to="/home" class="app-topbar__brand">
        <img class="app-topbar__logo app-logo" src="/nest.png" alt="ApiNest" />
        <span class="app-topbar__title">ApiNest</span>
      </router-link>
      <WorkspaceProjectSelector />
    </div>

    <div class="app-topbar__actions">
      <el-tooltip :content="themeTooltip" placement="bottom">
        <button type="button" class="app-topbar__action" @click="toggleTheme">
          <el-icon :size="20">
            <Sunny v-if="theme === 'dark'" />
            <Moon v-else />
          </el-icon>
        </button>
      </el-tooltip>

      <LocaleSwitcher />

      <el-tooltip v-if="user?.is_admin" :content="t('topbar.settings')" placement="bottom">
        <button type="button" class="app-topbar__action" @click="goSettings">
          <el-icon :size="20"><Setting /></el-icon>
        </button>
      </el-tooltip>

      <el-dropdown trigger="click" popper-class="app-user-menu" @command="handleCommand">
        <button type="button" class="app-topbar__action app-topbar__action--avatar">
          <el-avatar :size="32" :src="avatarSrc || undefined">{{ avatarText }}</el-avatar>
        </button>
        <template #dropdown>
          <div class="app-user-menu__inner">
            <div class="app-user-menu__profile">
              <el-avatar :size="42" :src="avatarSrc || undefined">{{ avatarText }}</el-avatar>
              <div class="app-user-menu__info">
                <span class="app-user-menu__name">{{ displayName }}</span>
                <span class="app-user-menu__account">{{ user?.account }}</span>
                <span v-if="user?.is_admin" class="app-user-menu__badge">{{ t('common.admin') }}</span>
              </div>
            </div>
            <div class="app-user-menu__divider" />
            <button type="button" class="app-user-menu__logout" @click="handleLogout">
              <el-icon :size="16"><SwitchButton /></el-icon>
              <span>{{ t('common.logout') }}</span>
            </button>
          </div>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<style scoped>
.app-topbar {
  position: sticky;
  top: 0;
  z-index: 200;
  height: 56px;
  width: 100%;
  padding: 0 20px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-topbar);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-topbar__left {
  display: flex;
  align-items: center;
  min-width: 0;
  flex: 1;
}

.app-topbar__brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--color-text);
  flex-shrink: 0;
}

.app-topbar__logo {
  width: 32px;
  height: 32px;
  object-fit: contain;
  display: block;
  flex-shrink: 0;
}

.app-topbar__title {
  font-size: 18px;
  font-weight: 700;
}

.app-topbar__actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.app-topbar__action {
  width: 36px;
  height: 36px;
  border: none;
  border-radius: 10px;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  transition: background-color 0.15s ease, color 0.15s ease;
}

.app-topbar__action:hover {
  background: var(--color-hover);
  color: var(--color-text);
}

.app-topbar__action--text {
  width: auto;
  min-width: 32px;
  padding: 0 8px;
  font-size: 13px;
  font-weight: 600;
}

.app-topbar__action--avatar {
  width: 32px;
  height: 32px;
}
</style>
