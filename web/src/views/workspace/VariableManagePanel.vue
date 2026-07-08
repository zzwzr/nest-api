<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  createEnvironmentVariable,
  deleteEnvironmentVariable,
  updateEnvironmentVariable,
} from '@/api/envvariable'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'

const { t } = useLocale()
const {
  activeTab,
  activeProject,
  activeWorkspaceId,
  activeProjectId,
  activeEnvironmentId,
  environmentVariables,
  loadingEnvironmentVariables,
  loadEnvironmentVariables,
} = useWorkspaceContext()

const editingId = ref<number | null>(null)
const deletingId = ref<number | null>(null)
const creating = ref(false)

const currentEnvironmentId = computed(() => activeTab.value?.environmentId ?? activeEnvironmentId.value)

const currentEnvironmentName = computed(() => activeTab.value?.label ?? '')

watch(
  () => currentEnvironmentId.value,
  (environmentId) => {
    if (!environmentId) return
    loadEnvironmentVariables()
  },
  { immediate: true },
)

async function handleCreate() {
  if (!activeWorkspaceId.value || !activeProjectId.value || !currentEnvironmentId.value) return

  try {
    const { value: key } = await ElMessageBox.prompt(
      t('environment.variableKeyPlaceholder'),
      t('environment.createVariable'),
      {
        inputPattern: /\S+/,
        inputErrorMessage: t('workspace.nameRequired'),
      },
    )
    if (!key?.trim()) return

    const { value } = await ElMessageBox.prompt(
      t('environment.variableValuePlaceholder'),
      t('environment.createVariable'),
      {
        inputValue: '',
      },
    )

    creating.value = true
    await createEnvironmentVariable(
      activeWorkspaceId.value,
      activeProjectId.value,
      currentEnvironmentId.value,
      {
        key: key.trim(),
        value: value ?? '',
      },
    )
    ElMessage.success(t('environment.variableCreateSuccess'))
    await loadEnvironmentVariables()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error instanceof Error ? error.message : t('environment.variableCreateFailed'))
    }
  } finally {
    creating.value = false
  }
}

async function handleEdit(variableId: number, currentKey: string, currentValue: string, description: string) {
  if (!activeWorkspaceId.value || !activeProjectId.value || !currentEnvironmentId.value) return

  try {
    const { value: key } = await ElMessageBox.prompt(
      t('environment.variableKeyPlaceholder'),
      t('environment.editVariable'),
      {
        inputValue: currentKey,
        inputPattern: /\S+/,
        inputErrorMessage: t('workspace.nameRequired'),
      },
    )
    if (!key?.trim()) return

    const { value } = await ElMessageBox.prompt(
      t('environment.variableValuePlaceholder'),
      t('environment.editVariable'),
      {
        inputValue: currentValue,
      },
    )

    editingId.value = variableId
    await updateEnvironmentVariable(
      activeWorkspaceId.value,
      activeProjectId.value,
      currentEnvironmentId.value,
      variableId,
      {
        key: key.trim(),
        value: value ?? '',
        description,
      },
    )
    ElMessage.success(t('environment.variableUpdateSuccess'))
    await loadEnvironmentVariables()
  } catch (error) {
    if (error !== 'cancel' && error !== 'close') {
      ElMessage.error(error instanceof Error ? error.message : t('environment.variableUpdateFailed'))
    }
  } finally {
    editingId.value = null
  }
}

async function handleDelete(variableId: number, key: string) {
  if (!activeWorkspaceId.value || !activeProjectId.value || !currentEnvironmentId.value) return

  try {
    await ElMessageBox.confirm(
      t('environment.variableDeleteConfirm', { name: key }),
      t('common.delete'),
      { type: 'warning' },
    )
  } catch {
    return
  }

  deletingId.value = variableId
  try {
    await deleteEnvironmentVariable(
      activeWorkspaceId.value,
      activeProjectId.value,
      currentEnvironmentId.value,
      variableId,
    )
    ElMessage.success(t('environment.variableDeleteSuccess'))
    await loadEnvironmentVariables()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('environment.variableDeleteFailed'))
  } finally {
    deletingId.value = null
  }
}
</script>

<template>
  <div class="variable-manage workspace-panel">
    <div class="variable-manage__toolbar workspace-panel__toolbar">
      <div>
        <h2>{{ t('environment.variables') }}</h2>
        <p v-if="activeProject && currentEnvironmentName">
          {{ t('environment.variableManageDesc', { project: activeProject.name, env: currentEnvironmentName }) }}
        </p>
      </div>
      <el-button
        type="primary"
        class="workspace-action-btn"
        :loading="creating"
        :disabled="!currentEnvironmentId"
        @click="handleCreate"
      >
        <span class="workspace-action-btn__plus">+</span>
        <span>{{ t('environment.createVariable') }}</span>
      </el-button>
    </div>

    <el-table
      v-loading="loadingEnvironmentVariables"
      :data="environmentVariables"
      class="workspace-data-table"
    >
      <el-table-column prop="key" :label="t('environment.variableKey')" min-width="160" />
      <el-table-column prop="value" :label="t('environment.variableValue')" min-width="200" show-overflow-tooltip />
      <el-table-column prop="description" :label="t('environment.variableDescription')" min-width="160" show-overflow-tooltip />
      <el-table-column prop="updated_at" :label="t('workspace.columns.updatedAt')" min-width="180" />
      <el-table-column :label="t('member.columns.actions')" width="160" fixed="right">
        <template #default="{ row }">
          <el-button
            type="primary"
            link
            :loading="editingId === row.id"
            @click="handleEdit(row.id, row.key, row.value, row.description)"
          >
            {{ t('common.edit') }}
          </el-button>
          <el-button
            type="danger"
            link
            :loading="deletingId === row.id"
            @click="handleDelete(row.id, row.key)"
          >
            {{ t('common.delete') }}
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<style scoped>
.variable-manage__toolbar {
  padding: 10px 10px 10px 8px;
}
</style>
