<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import InterfaceStatusRadio from '@/components/interface/InterfaceStatusRadio.vue'
import InterfaceUrlBar from '@/components/interface/InterfaceUrlBar.vue'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'
import type { ApiTreeNode, HttpMethod, HttpProtocol, InterfaceStatus } from '@/types/workspace'

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
  protocol: 'HTTP' as HttpProtocol,
  method: 'POST' as HttpMethod,
  url: '/',
  folderId: null as number | null,
  name: '',
  status: 1 as InterfaceStatus,
})

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
  () => [form.protocol, form.method, form.url, form.folderId, form.name, form.status],
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
    form.protocol = 'HTTP'
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
  <div class="interface-create interface-editor">
    <div class="interface-create__toolbar">
      <el-button type="primary" :loading="saving" @click="handleSave">
        {{ t('common.save') }}
      </el-button>
      <el-button @click="handleCancel">{{ t('common.cancel') }}</el-button>
    </div>

    <div class="interface-create__section">
      <div class="interface-create__row">
        <InterfaceUrlBar
          v-model:protocol="form.protocol"
          v-model:method="form.method"
          v-model:url="form.url"
        />
      </div>

      <div class="interface-create__row interface-meta-row">
        <el-select
          v-model="form.folderId"
          class="interface-meta-row__folder"
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
          class="interface-meta-row__name"
          :placeholder="t('workspace.apiNamePlaceholder')"
        />
      </div>

      <div class="interface-create__status-row">
        <span class="interface-create__label">{{ t('workspace.columns.status') }}</span>
        <InterfaceStatusRadio v-model="form.status" />
      </div>
    </div>
  </div>
</template>

<style scoped>
.interface-create {
  padding: 12px 0 32px;
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
  gap: 10px;
  padding-bottom: 20px;
  padding-left: var(--color-interface-section-padding-x);
  padding-right: var(--color-interface-section-padding-x);
  border-bottom: 1px solid var(--color-border);
}

.interface-create__row {
  display: flex;
  align-items: center;
  width: 100%;
}

.interface-create__status-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.interface-create__label {
  color: var(--color-interface-field-text);
  font-size: 14px;
  flex-shrink: 0;
}
</style>
