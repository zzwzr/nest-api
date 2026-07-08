<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createEnvironment,
  deleteEnvironment,
  updateEnvironment,
} from '@/api/environment'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'

const { t } = useLocale()
const {
  activeProject,
  activeWorkspaceId,
  activeProjectId,
  environments,
  loadingEnvironments,
  refreshEnvironments,
} = useWorkspaceContext()

const editingId = ref<number | null>(null)
const deletingId = ref<number | null>(null)
const creating = ref(false)

async function handleCreate() {
  if (!activeWorkspaceId.value || !activeProjectId.value) return

  try {
    const { value } = await ElMessageBox.prompt(
      t('environment.namePlaceholder'),
      t('environment.create'),
      {
        inputPattern: /\S+/,
        inputErrorMessage: t('workspace.nameRequired'),
      },
    )
    if (!value?.trim()) return

    creating.value = true
    await createEnvironment(activeWorkspaceId.value, activeProjectId.value, {
      name: value.trim(),
    })
    ElMessage.success(t('environment.createSuccess'))
    await refreshEnvironments()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error instanceof Error ? error.message : t('environment.createFailed'))
    }
  } finally {
    creating.value = false
  }
}

async function handleRename(environmentId: number, currentName: string, baseUrl: string) {
  if (!activeWorkspaceId.value || !activeProjectId.value) return

  try {
    const { value } = await ElMessageBox.prompt(
      t('environment.namePlaceholder'),
      t('environment.rename'),
      {
        inputValue: currentName,
        inputPattern: /\S+/,
        inputErrorMessage: t('workspace.nameRequired'),
      },
    )
    if (!value?.trim()) return

    editingId.value = environmentId
    await updateEnvironment(activeWorkspaceId.value, activeProjectId.value, environmentId, {
      name: value.trim(),
      base_url: baseUrl,
    })
    ElMessage.success(t('environment.updateSuccess'))
    await refreshEnvironments()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error instanceof Error ? error.message : t('environment.updateFailed'))
    }
  } finally {
    editingId.value = null
  }
}

async function handleDelete(environmentId: number, name: string) {
  if (!activeWorkspaceId.value || !activeProjectId.value) return

  try {
    await ElMessageBox.confirm(
      t('environment.deleteConfirm', { name }),
      t('environment.delete'),
      { type: 'warning' },
    )
  } catch {
    return
  }

  deletingId.value = environmentId
  try {
    await deleteEnvironment(activeWorkspaceId.value, activeProjectId.value, environmentId)
    ElMessage.success(t('environment.deleteSuccess'))
    await refreshEnvironments()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('environment.deleteFailed'))
  } finally {
    deletingId.value = null
  }
}

async function handleSetDefault(environmentId: number, name: string, baseUrl: string) {
  if (!activeWorkspaceId.value || !activeProjectId.value) return

  try {
    await updateEnvironment(activeWorkspaceId.value, activeProjectId.value, environmentId, {
      name,
      base_url: baseUrl,
      is_default: true,
    })
    ElMessage.success(t('environment.setDefaultSuccess'))
    await refreshEnvironments()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('environment.updateFailed'))
  }
}

watch(activeProjectId, () => {
  refreshEnvironments()
})
</script>

<template>
  <div class="environment-manage workspace-panel">
    <div class="environment-manage__toolbar workspace-panel__toolbar">
      <div>
        <h2>{{ t('workspace.modules.environment') }}</h2>
        <p v-if="activeProject">
          {{ t('environment.manageDesc', { name: activeProject.name }) }}
        </p>
      </div>
      <el-button type="primary" class="workspace-action-btn" :loading="creating" @click="handleCreate">
        <span>{{ t('environment.create') }}</span>
      </el-button>
    </div>

    <el-table v-loading="loadingEnvironments" :data="environments" class="workspace-data-table">
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" :label="t('workspace.columns.name')" min-width="140" />
      <el-table-column prop="base_url" :label="t('environment.baseUrl')" min-width="200" />
      <el-table-column :label="t('environment.default')" width="100">
        <template #default="{ row }">
          <span v-if="row.is_default" class="environment-manage__default">{{ t('environment.yes') }}</span>
          <el-button
            v-else
            type="primary"
            link
            @click="handleSetDefault(row.id, row.name, row.base_url)"
          >
            {{ t('environment.setDefault') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" :label="t('workspace.columns.createdAt')" min-width="180" />
      <el-table-column :label="t('member.columns.actions')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            link
            :loading="editingId === row.id"
            @click="handleRename(row.id, row.name, row.base_url)"
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
  </div>
</template>

<style scoped>
.environment-manage__toolbar {
  padding: 10px 10px 10px 8px;
}

.environment-manage__default {
  color: #49cc90;
  font-size: 14px;
}
</style>
