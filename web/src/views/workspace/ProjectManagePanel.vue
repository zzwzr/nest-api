<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteProject, updateProject } from '@/api/project'
import MemberManagePanel from '@/views/workspace/MemberManagePanel.vue'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'

type ProjectTab = 'projects' | 'members'

const { t } = useLocale()
const {
  projects,
  activeWorkspace,
  activeWorkspaceId,
  loadingProjects,
  refreshProjects,
  openCreateProject,
} = useWorkspaceContext()

const activeTab = ref<ProjectTab>('projects')
const editingId = ref<number | null>(null)
const deletingId = ref<number | null>(null)

const tabs = computed(() => [
  { key: 'projects' as ProjectTab, label: t('workspace.projectManage') },
  { key: 'members' as ProjectTab, label: t('member.title') },
])

async function handleRename(projectId: number, currentName: string) {
  try {
    const { value } = await ElMessageBox.prompt(
      t('workspace.renameProjectPrompt'),
      t('workspace.renameProject'),
      {
        inputValue: currentName,
        inputPattern: /\S+/,
        inputErrorMessage: t('workspace.nameRequired'),
      },
    )
    if (!activeWorkspaceId.value || !value?.trim()) return

    editingId.value = projectId
    await updateProject(activeWorkspaceId.value, projectId, value.trim())
    ElMessage.success(t('workspace.updateProjectSuccess'))
    await refreshProjects()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error instanceof Error ? error.message : t('workspace.updateProjectFailed'))
    }
  } finally {
    editingId.value = null
  }
}

async function handleDelete(projectId: number, name: string) {
  if (!activeWorkspaceId.value) return

  try {
    await ElMessageBox.confirm(
      t('workspace.deleteProjectConfirm', { name }),
      t('workspace.deleteProject'),
      { type: 'warning' },
    )
  } catch {
    return
  }

  deletingId.value = projectId
  try {
    await deleteProject(activeWorkspaceId.value, projectId)
    ElMessage.success(t('workspace.deleteProjectSuccess'))
    await refreshProjects()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('workspace.deleteProjectFailed'))
  } finally {
    deletingId.value = null
  }
}

watch(activeWorkspaceId, () => {
  activeTab.value = 'projects'
})
</script>

<template>
  <div class="project-manage">
    <div class="project-manage__tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="project-manage__tab"
        :class="{ 'project-manage__tab--active': activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <template v-if="activeTab === 'projects'">
      <div class="project-manage__toolbar">
        <div>
          <h2>{{ t('workspace.projectManage') }}</h2>
          <p v-if="activeWorkspace">
            {{ t('workspace.projectManageDesc', { name: activeWorkspace.name }) }}
          </p>
        </div>
        <el-button type="primary" @click="openCreateProject">
          <el-icon><Plus /></el-icon>
          {{ t('workspace.createProject') }}
        </el-button>
      </div>

      <el-table v-loading="loadingProjects" :data="projects" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="t('workspace.columns.name')" min-width="160" />
        <el-table-column prop="creator_name" :label="t('workspace.columns.creator')" min-width="120" />
        <el-table-column prop="created_at" :label="t('workspace.columns.createdAt')" min-width="180" />
        <el-table-column :label="t('member.columns.actions')" width="160" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              :loading="editingId === row.id"
              @click="handleRename(row.id, row.name)"
            >
              {{ t('common.rename') }}
            </el-button>
            <el-button
              type="danger"
              link
              :loading="deletingId === row.id"
              @click="handleDelete(row.id, row.name)"
            >
              {{ t('common.delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </template>

    <MemberManagePanel v-else />
  </div>
</template>

<style scoped>
.project-manage {
  min-height: 100%;
}

.project-manage__tabs {
  display: flex;
  gap: 4px;
  padding: 16px 20px 0;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-surface);
}

.project-manage__tab {
  height: 40px;
  padding: 0 16px;
  border: none;
  border-bottom: 2px solid transparent;
  background: transparent;
  color: var(--color-text-secondary);
  font-size: 15px;
  cursor: pointer;
  transition: color 0.15s ease, border-color 0.15s ease;
}

.project-manage__tab:hover {
  color: var(--color-text);
}

.project-manage__tab--active {
  color: var(--color-primary-light);
  border-bottom-color: var(--color-primary);
  font-weight: 600;
}

.project-manage__toolbar {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  padding: 20px;
}

.project-manage__toolbar h2 {
  margin: 0 0 6px;
  font-size: 20px;
  font-weight: 700;
  color: var(--color-text);
}

.project-manage__toolbar p {
  margin: 0;
  font-size: 14px;
  color: var(--color-text-secondary);
}
</style>
