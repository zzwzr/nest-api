<script setup lang="ts">
import { computed } from 'vue'
import { Close } from '@element-plus/icons-vue'
import ProjectManagePanel from '@/views/workspace/ProjectManagePanel.vue'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'
import type { AppModule } from '@/types/workspace'

const { t } = useLocale()
const { activeModule, activeProject, selectedApiId, contextMode } = useWorkspaceContext()

const moduleTitleMap: Record<AppModule, string> = {
  api: 'workspace.modules.api',
  'quick-test': 'workspace.modules.quickTest',
  environment: 'workspace.modules.environment',
  project: 'workspace.modules.project',
}

const pageTitle = computed(() => {
  if (activeModule.value === 'api' && selectedApiId.value) {
    return t('workspace.apiDetail')
  }
  if (activeModule.value === 'api' && activeProject.value) {
    return t('workspace.apiList')
  }
  return t(moduleTitleMap[activeModule.value])
})

const pageDesc = computed(() => {
  if (activeModule.value === 'api') {
    if (!activeProject.value) return t('workspace.apiDescNoProject')
    if (!selectedApiId.value) return t('workspace.apiDescList')
    return t('workspace.apiDescDetail')
  }
  if (activeModule.value === 'quick-test') return t('workspace.quickTestDesc')
  if (activeModule.value === 'environment') return t('workspace.environmentDesc')
  return t('workspace.projectDesc')
})

const showApiTable = computed(
  () => activeModule.value === 'api' && contextMode.value === 'project' && !selectedApiId.value,
)

const showProjectManage = computed(() => activeModule.value === 'project')

const mockApiRows = [
  { id: 1, name: '退出登录', status: 'published', method: 'POST', url: '/api/auth/logout', group: '账号相关' },
  { id: 2, name: '获取用户信息', status: 'published', method: 'GET', url: '/api/user/info', group: '账号相关' },
  { id: 3, name: '修改密码', status: 'testing', method: 'PUT', url: '/api/user/password', group: '用户模块' },
  { id: 4, name: '创建订单', status: 'published', method: 'POST', url: '/api/order/create', group: '订单模块' },
]
</script>

<template>
  <div class="workspace-main">
    <div class="workspace-main__tabs">
      <div class="workspace-main__tab workspace-main__tab--active">
        <span>{{ pageTitle }}</span>
        <button type="button" class="workspace-main__tab-close">
          <el-icon :size="12"><Close /></el-icon>
        </button>
      </div>
    </div>

    <div v-if="!showProjectManage" class="workspace-main__toolbar">
      <nav class="workspace-main__subnav">
        <a class="workspace-main__subnav-item workspace-main__subnav-item--active" href="#">
          {{ t('workspace.subnav.api') }}
        </a>
        <a v-if="activeModule === 'api'" class="workspace-main__subnav-item" href="#">
          {{ t('workspace.subnav.group') }}
        </a>
      </nav>

      <div v-if="activeModule === 'api' && activeProject" class="workspace-main__actions">
        <button type="button" class="workspace-main__btn workspace-main__btn--primary">
          + {{ t('workspace.addApi') }}
        </button>
      </div>
    </div>

    <div class="workspace-main__content">
      <template v-if="showProjectManage">
        <ProjectManagePanel />
      </template>

      <template v-else-if="showApiTable">
        <el-table :data="mockApiRows" stripe class="workspace-main__table">
          <el-table-column :label="t('workspace.columns.id')" prop="id" width="60" />
          <el-table-column :label="t('workspace.columns.name')" prop="name" min-width="140" />
          <el-table-column :label="t('workspace.columns.status')" width="100">
            <template #default="{ row }">
              <span
                class="workspace-main__status"
                :class="`workspace-main__status--${row.status}`"
              >
                {{ row.status === 'published' ? t('workspace.status.published') : t('workspace.status.testing') }}
              </span>
            </template>
          </el-table-column>
          <el-table-column :label="t('workspace.columns.method')" width="90">
            <template #default="{ row }">
              <span class="workspace-main__method" :class="`workspace-main__method--${row.method.toLowerCase()}`">
                {{ row.method }}
              </span>
            </template>
          </el-table-column>
          <el-table-column :label="t('workspace.columns.url')" prop="url" min-width="200" />
          <el-table-column :label="t('workspace.columns.group')" prop="group" width="120" />
        </el-table>
        <div class="workspace-main__footer">
          {{ t('workspace.loadedRecords', { count: String(mockApiRows.length) }) }}
        </div>
      </template>

      <template v-else>
        <div class="workspace-main__placeholder">
          <h2>{{ pageTitle }}</h2>
          <p>{{ pageDesc }}</p>
        </div>
      </template>
    </div>
  </div>
</template>

<style scoped>
.workspace-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  background: var(--color-bg);
  overflow: hidden;
}

.workspace-main__tabs {
  display: flex;
  align-items: flex-end;
  gap: 2px;
  padding: 0 16px;
  height: 42px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-surface);
}

.workspace-main__tab {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  height: 34px;
  padding: 0 14px;
  border-radius: 8px 8px 0 0;
  font-size: 14px;
  color: var(--color-text-secondary);
  background: transparent;
  border: 1px solid transparent;
  cursor: pointer;
}

.workspace-main__tab--active {
  color: var(--color-text);
  background: var(--color-bg);
  border-color: var(--color-border);
  border-bottom-color: var(--color-bg);
  margin-bottom: -1px;
}

.workspace-main__tab-close {
  width: 16px;
  height: 16px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0;
}

.workspace-main__tab-close:hover {
  background: var(--color-hover);
  color: var(--color-text);
}

.workspace-main__toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 12px 20px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-surface);
}

.workspace-main__subnav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.workspace-main__subnav-item {
  font-size: 15px;
  color: var(--color-text-secondary);
  text-decoration: none;
  padding-bottom: 2px;
  border-bottom: 2px solid transparent;
  transition: color 0.15s ease, border-color 0.15s ease;
}

.workspace-main__subnav-item:hover {
  color: var(--color-text);
}

.workspace-main__subnav-item--active {
  color: var(--color-primary-light);
  border-bottom-color: var(--color-primary);
  font-weight: 500;
}

.workspace-main__actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.workspace-main__btn {
  height: 36px;
  padding: 0 16px;
  border: 1px solid var(--color-border);
  border-radius: 8px;
  background: var(--color-surface);
  color: var(--color-text);
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.workspace-main__btn:hover {
  background: var(--color-hover);
}

.workspace-main__btn--primary {
  background: var(--color-primary);
  border-color: var(--color-primary);
  color: #fff;
}

.workspace-main__btn--primary:hover {
  background: var(--color-primary-light);
  border-color: var(--color-primary-light);
}

.workspace-main__content {
  flex: 1;
  overflow: auto;
  padding: 0;
}

.workspace-main__table {
  width: 100%;
}

.workspace-main__footer {
  padding: 12px 20px;
  font-size: 14px;
  color: var(--color-text-secondary);
  border-top: 1px solid var(--color-border);
  background: var(--color-surface);
}

.workspace-main__status {
  display: inline-block;
  padding: 3px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 500;
}

.workspace-main__status--published {
  color: #49cc90;
  background: rgba(73, 204, 144, 0.12);
}

.workspace-main__status--testing {
  color: #fca130;
  background: rgba(252, 161, 48, 0.12);
}

.workspace-main__method {
  display: inline-block;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 700;
  font-family: 'Consolas', 'Monaco', monospace;
}

.workspace-main__method--get {
  color: #61affe;
  background: rgba(97, 175, 254, 0.12);
}

.workspace-main__method--post {
  color: #49cc90;
  background: rgba(73, 204, 144, 0.12);
}

.workspace-main__method--put {
  color: #fca130;
  background: rgba(252, 161, 48, 0.12);
}

.workspace-main__placeholder {
  padding: 48px 32px;
  text-align: center;
}

.workspace-main__placeholder h2 {
  margin: 0 0 12px;
  font-size: 22px;
  font-weight: 600;
  color: var(--color-text);
}

.workspace-main__placeholder p {
  margin: 0;
  font-size: 15px;
  color: var(--color-text-secondary);
  line-height: 1.6;
}
</style>
