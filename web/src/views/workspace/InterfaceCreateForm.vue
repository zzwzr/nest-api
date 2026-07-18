<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import InterfaceEditorPanels from '@/interface/InterfaceEditorPanels.vue'
import InterfaceStatusRadio from '@/interface/InterfaceStatusRadio.vue'
import InterfaceUrlBar from '@/interface/InterfaceUrlBar.vue'
import {
  useInterfaceEditorDirty,
  useSaveShortcut,
} from '@/composables/useInterfaceEditorDirty'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'
import {
  buildInterfaceSavePayload,
  createEmptyInterfaceEditorForm,
  defaultResponseExample,
  defaultResponseResult,
  interfaceEditorSnapshot,
} from '@/utils/interface-editor-form'
import { collectFolderOptions } from '@/utils/interface-folder-options'

const { t } = useLocale()
const emit = defineEmits<{
  'request-close': []
}>()

const {
  apiTree,
  selectedFolder,
  submitCreateInterface,
  parseFolderId,
  setActiveTabDirty,
} = useWorkspaceContext()

const saving = ref(false)
const form = reactive(createEmptyInterfaceEditorForm({ method: 'POST', status: 1 }))

const folderOptions = computed(() => collectFolderOptions(apiTree.value, parseFolderId))

const { beginSuppress, captureBaseline } = useInterfaceEditorDirty({
  getSnapshot: () => interfaceEditorSnapshot(form),
  setDirty: setActiveTabDirty,
  enabled: () => !saving.value,
  watchSource: () => form,
})

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

onMounted(async () => {
  beginSuppress()
  form.responseResults = [defaultResponseResult(t('workspace.interfaceForm.defaultSuccess'))]
  form.responseExamples = [defaultResponseExample(t('workspace.interfaceForm.successExample'))]
  if (!form.folderId && folderOptions.value.length > 0) {
    form.folderId = folderOptions.value[0].folderId
  }
  await captureBaseline()
})

async function handleSave(): Promise<boolean> {
  if (saving.value) return false
  const folderId = form.folderId
  if (!folderId) {
    ElMessage.warning(t('workspace.folderRequired'))
    return false
  }
  if (!form.name.trim()) {
    ElMessage.warning(t('workspace.nameRequired'))
    return false
  }

  saving.value = true
  try {
    await submitCreateInterface({
      ...buildInterfaceSavePayload(form),
      folder_id: folderId,
    })
    ElMessage.success(t('workspace.createApiSuccess'))
    return true
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('workspace.createApiFailed'))
    return false
  } finally {
    saving.value = false
  }
}

useSaveShortcut(handleSave, () => !saving.value)

defineExpose({
  save: handleSave,
})

function emitClose() {
  emit('request-close')
}
</script>

<template>
  <div class="interface-create interface-editor">
    <div class="interface-create__toolbar">
      <el-button type="primary" :loading="saving" @click="handleSave">
        {{ t('common.save') }}
      </el-button>
      <el-button @click="emitClose">{{ t('common.cancel') }}</el-button>
    </div>

    <div class="interface-create__scroll">
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

      <InterfaceEditorPanels
        v-model:request-headers="form.requestHeaders"
        v-model:query-params="form.queryParams"
        v-model:request-body="form.requestBody"
        v-model:response-headers="form.responseHeaders"
        v-model:response-results="form.responseResults"
        v-model:response-examples="form.responseExamples"
      />
    </div>
  </div>
</template>

<style scoped>
.interface-create {
  display: flex;
  flex-direction: column;
  flex: 1;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.interface-create__toolbar {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-workspace-content);
}

.interface-create__scroll {
  flex: 1;
  min-height: 0;
  overflow: auto;
  overscroll-behavior: contain;
  padding: 12px 0 32px;
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
