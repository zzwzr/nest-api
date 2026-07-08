<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'
import type { ApiTreeNode, HttpMethod } from '@/types/workspace'

const { t } = useLocale()
const {
  apiTree,
  selectedFolder,
  submitCreateInterface,
  closeCreateApi,
  parseFolderId,
  markActiveTabDirty,
} = useWorkspaceContext()

const saving = ref(false)

const form = reactive({
  method: 'POST' as HttpMethod,
  url: '/',
  folderId: null as number | null,
  name: '',
  status: 1 as 1 | 2,
})

const methodOptions: HttpMethod[] = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']

const statusOptions = computed(() => [
  { value: 1, label: t('workspace.status.published') },
  { value: 2, label: t('workspace.status.testing') },
])

interface FolderOption {
  folderId: number
  label: string
}

function collectFolderOptions(nodes: ApiTreeNode[], depth = 0): FolderOption[] {
  const result: FolderOption[] = []
  for (const node of nodes) {
    if (node.type !== 'folder') continue
    const folderId = parseFolderId(node.id)
    if (folderId) {
      result.push({
        folderId,
        label: `${'\u3000'.repeat(depth)}${node.name}`,
      })
    }
    if (node.children?.length) {
      result.push(...collectFolderOptions(node.children, depth + 1))
    }
  }
  return result
}

const folderOptions = computed(() => collectFolderOptions(apiTree.value))

watch(
  () => [form.method, form.url, form.folderId, form.name, form.status],
  () => {
    markActiveTabDirty()
  },
)

watch(
  selectedFolder,
  (folder) => {
    if (!folder) return
    form.folderId = parseFolderId(folder.id)
  },
  { immediate: true },
)

watch(folderOptions, (options) => {
  if (!form.folderId && options.length > 0) {
    form.folderId = options[0].folderId
  }
})

async function handleSave() {
  if (!form.folderId) {
    ElMessage.warning(t('workspace.folderRequired'))
    return
  }
  if (!form.name.trim()) {
    ElMessage.warning(t('workspace.nameRequired'))
    return
  }

  saving.value = true
  try {
    await submitCreateInterface(
      form.folderId,
      form.name.trim(),
      form.method,
      form.url.trim(),
      form.status,
    )
    ElMessage.success(t('workspace.createApiSuccess'))
    form.name = ''
    form.url = '/'
    form.method = 'POST'
    form.status = 1
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('workspace.createApiFailed'))
  } finally {
    saving.value = false
  }
}

function handleCancel() {
  closeCreateApi()
}
</script>

<template>
  <div class="interface-create">
    <div class="interface-create__toolbar">
      <el-button type="primary" :loading="saving" @click="handleSave">
        {{ t('common.save') }}
      </el-button>
      <el-button @click="handleCancel">{{ t('common.cancel') }}</el-button>
    </div>

    <div class="interface-create__section">
      <div class="interface-create__row">
        <el-select
          v-model="form.method"
          class="interface-create__method"
          size="default"
          popper-class="app-action-dropdown"
        >
          <el-option v-for="item in methodOptions" :key="item" :label="item" :value="item" />
        </el-select>
        <el-select model-value="HTTP" class="interface-create__protocol" size="default" disabled>
          <el-option label="HTTP" value="HTTP" />
        </el-select>
        <el-input
          v-model="form.url"
          class="interface-create__url"
          :placeholder="t('workspace.urlPlaceholder')"
        />
      </div>

      <div class="interface-create__row">
        <el-select
          v-model="form.folderId"
          class="interface-create__folder"
          :placeholder="t('workspace.selectFolder')"
          size="default"
          popper-class="app-action-dropdown"
        >
          <el-option
            v-for="item in folderOptions"
            :key="item.folderId"
            :label="item.label"
            :value="item.folderId"
          />
        </el-select>
        <el-input
          v-model="form.name"
          class="interface-create__name"
          :placeholder="t('workspace.apiNamePlaceholder')"
        />
      </div>

      <div class="interface-create__status-row">
        <span class="interface-create__label">{{ t('workspace.columns.status') }}</span>
        <el-radio-group v-model="form.status" class="interface-create__status-group">
          <el-radio v-for="item in statusOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </el-radio>
        </el-radio-group>
      </div>
    </div>

    <div class="interface-create__panels">
      <details class="interface-create__panel" open>
        <summary>{{ t('workspace.interfaceForm.moreSettings') }}</summary>
        <p class="interface-create__panel-desc">{{ t('workspace.comingSoon') }}</p>
      </details>
      <details class="interface-create__panel">
        <summary>{{ t('workspace.interfaceForm.description') }}</summary>
        <p class="interface-create__panel-desc">{{ t('workspace.comingSoon') }}</p>
      </details>
      <details class="interface-create__panel">
        <summary>{{ t('workspace.interfaceForm.requestParams') }}</summary>
        <p class="interface-create__panel-desc">{{ t('workspace.comingSoon') }}</p>
      </details>
      <details class="interface-create__panel">
        <summary>{{ t('workspace.interfaceForm.response') }}</summary>
        <p class="interface-create__panel-desc">{{ t('workspace.comingSoon') }}</p>
      </details>
    </div>
  </div>
</template>

<style scoped>
.interface-create {
  padding: 12px 16px 32px;
  background: #212121;
  color: #e1e1e1;
  font-size: 14px;
}

.interface-create :deep(.el-input__inner),
.interface-create :deep(.el-input__wrapper),
.interface-create :deep(.el-select__wrapper),
.interface-create :deep(.el-select__selected-item),
.interface-create :deep(.el-radio__label),
.interface-create :deep(.el-button:not(.el-button--primary)) {
  color: #e1e1e1;
}

.interface-create :deep(.el-input__inner::placeholder) {
  color: #8a8a8a;
}

.interface-create__toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}

.interface-create__section {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding-bottom: 20px;
  border-bottom: 1px solid var(--color-border);
}

.interface-create__row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.interface-create__method {
  width: 110px;
  flex-shrink: 0;
}

.interface-create__protocol {
  width: 100px;
  flex-shrink: 0;
}

.interface-create__url {
  flex: 1;
  min-width: 0;
}

.interface-create__folder {
  width: 220px;
  flex-shrink: 0;
}

.interface-create__name {
  flex: 1;
  min-width: 0;
}

.interface-create__status-row {
  display: flex;
  align-items: center;
  gap: 16px;
  flex-wrap: wrap;
}

.interface-create__label {
  color: #e1e1e1;
  font-size: 14px;
  flex-shrink: 0;
}

.interface-create__status-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 16px;
}

.interface-create__panels {
  padding-top: 8px;
}

.interface-create__panel {
  border-bottom: 1px solid var(--color-border);
}

.interface-create__panel summary {
  padding: 14px 0;
  font-size: 14px;
  font-weight: 500;
  color: #e1e1e1;
  cursor: pointer;
  list-style: none;
}

.interface-create__panel summary::-webkit-details-marker {
  display: none;
}

.interface-create__panel-desc {
  margin: 0 0 16px;
  font-size: 14px;
  color: #e1e1e1;
}
</style>
