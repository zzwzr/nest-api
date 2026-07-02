<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { Moon, Setting, Sunny } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useAuth } from '@/composables/useAuth'
import { useLocale, type LocaleCode } from '@/composables/useLocale'
import { useTheme } from '@/composables/useTheme'

const router = useRouter()
const { user, logout } = useAuth()
const { theme, toggleTheme } = useTheme()
const { locale, t, setLocale } = useLocale()

const avatarText = computed(() => {
  const name = user.value?.name || user.value?.account || '?'
  return name.slice(0, 1).toUpperCase()
})

const avatarSrc = computed(() => user.value?.avatar || '')

const themeTooltip = computed(() =>
  theme.value === 'light' ? t('topbar.themeLight') : t('topbar.themeDark'),
)

const localeLabel = computed(() => (locale.value === 'zh-CN' ? '中文' : 'EN'))

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

function handleLocaleChange(code: LocaleCode) {
  setLocale(code)
}
</script>

<template>
  <header class="app-topbar">
    <router-link to="/home" class="app-topbar__brand">
      <span class="app-topbar__logo">N</span>
      <span class="app-topbar__title">ApiNest</span>
    </router-link>

    <div class="app-topbar__actions">
      <el-tooltip :content="themeTooltip" placement="bottom">
        <button type="button" class="app-topbar__action" @click="toggleTheme">
          <el-icon :size="18">
            <Sunny v-if="theme === 'dark'" />
            <Moon v-else />
          </el-icon>
        </button>
      </el-tooltip>

      <el-dropdown trigger="click" @command="handleLocaleChange">
        <button type="button" class="app-topbar__action app-topbar__action--text">
          {{ localeLabel }}
        </button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh-CN">中文</el-dropdown-item>
            <el-dropdown-item command="en-US">English</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>

      <el-tooltip v-if="user?.is_admin" :content="t('topbar.settings')" placement="bottom">
        <button type="button" class="app-topbar__action" @click="goSettings">
          <el-icon :size="18"><Setting /></el-icon>
        </button>
      </el-tooltip>

      <el-dropdown trigger="click" @command="handleCommand">
        <button type="button" class="app-topbar__action app-topbar__action--avatar">
          <el-avatar :size="32" :src="avatarSrc || undefined">{{ avatarText }}</el-avatar>
        </button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item disabled>
              <div class="app-topbar__profile">
                <strong>{{ user?.name || user?.account }}</strong>
                <span>{{ user?.account }}</span>
              </div>
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">{{ t('common.logout') }}</el-dropdown-item>
          </el-dropdown-menu>
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
  height: 50px;
  width: 100%;
  padding: 0 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-topbar);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.app-topbar__brand {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  color: var(--color-text);
  flex-shrink: 0;
}

.app-topbar__logo {
  width: 28px;
  height: 28px;
  border-radius: 8px;
  background: linear-gradient(135deg, var(--color-primary), var(--color-primary-light));
  color: #fff;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
}

.app-topbar__title {
  font-size: 16px;
  font-weight: 700;
}

.app-topbar__actions {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.app-topbar__action {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
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
  background: rgba(255, 255, 255, 0.06);
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

.app-topbar__profile {
  display: flex;
  flex-direction: column;
  gap: 2px;
  line-height: 1.3;
}

.app-topbar__profile span {
  font-size: 12px;
  color: var(--color-text-secondary);
}
</style>
