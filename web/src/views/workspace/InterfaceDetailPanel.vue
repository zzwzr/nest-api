<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { fetchInterfaceDetail } from '@/api/interface'
import InterfaceRequestParams from '@/components/interface/InterfaceRequestParams.vue'
import InterfaceResponseHeaderTable from '@/components/interface/InterfaceResponseHeaderTable.vue'
import InterfaceResponseExamples from '@/components/interface/InterfaceResponseExamples.vue'
import InterfaceResponseResults from '@/components/interface/InterfaceResponseResults.vue'
import InterfaceStatusRadio from '@/components/interface/InterfaceStatusRadio.vue'
import InterfaceUrlBar from '@/components/interface/InterfaceUrlBar.vue'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext, parseApiId } from '@/composables/useWorkspaceContext'
import {
  compactFieldTree,
  compactResponseFieldTree,
  emptyFieldNode,
  fieldTreeFromApi,
  responseFieldTreeFromApi,
} from '@/utils/interface-field-tree'
import {
  enrichResponseExample,
  normalizeResponseExampleForSave,
} from '@/utils/response-example-format'
import {
  compactParamRows,
  emptyParamRow,
  ensureTrailingEmptyRow,
  type ParamRow,
} from '@/utils/interface-params'
import type {
  ApiTreeNode,
  HttpMethod,
  HttpProtocol,
  InterfaceDetail,
  InterfaceRequestBody,
  InterfaceResponseExample,
  InterfaceResponseResult,
  InterfaceStatus,
} from '@/types/workspace'

type DetailMode = 'doc' | 'edit'

const { t } = useLocale()
const {
  activeModuleTab,
  activeWorkspaceId,
  activeProjectId,
  apiTree,
  parseFolderId,
  submitUpdateInterface,
  markActiveTabDirty,
} = useWorkspaceContext()

const mode = ref<DetailMode>('edit')
const loading = ref(false)
const saving = ref(false)

const openPanels = reactive({
  requestParams: true,
  responseHeader: true,
  responseResult: true,
  responseExample: true,
})

const form = reactive({
  protocol: 'HTTP' as HttpProtocol,
  method: 'GET' as HttpMethod,
  url: '/',
  folderId: null as number | null,
  name: '',
  status: 2 as InterfaceStatus,
  requestHeaders: [emptyParamRow()] as ParamRow[],
  requestBody: {
    format: 'json',
    data_type: 'Object',
    raw: '',
    fields: [emptyFieldNode()],
  } as InterfaceRequestBody,
  queryParams: [emptyParamRow()] as ParamRow[],
  responseHeaders: [emptyParamRow()] as ParamRow[],
  responseResults: [] as InterfaceResponseResult[],
  responseExamples: [] as InterfaceResponseExample[],
})

const interfaceId = computed(() => {
  const tab = activeModuleTab.value
  if (!tab?.apiId) return null
  return parseApiId(tab.apiId)
})

const readOnly = computed(() => mode.value === 'doc')

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

function toParamRows(items: { name: string; type: string; required: boolean; description: string; example: string }[]): ParamRow[] {
  const rows = items.map((item) => ({
    name: item.name,
    type: item.type || 'string',
    required: item.required,
    description: item.description,
    example: item.example,
  }))
  return ensureTrailingEmptyRow(rows.length ? rows : [emptyParamRow()])
}

function defaultResult(): InterfaceResponseResult {
  return {
    name: t('workspace.interfaceForm.defaultSuccess'),
    status_code: 200,
    format: 'JSON',
    data_type: 'Object',
    fields: responseFieldTreeFromApi([]),
  }
}

function defaultExample(): InterfaceResponseExample {
  return normalizeResponseExampleForSave({
    name: t('workspace.interfaceForm.successExample'),
    status_code: 200,
    content_type: 'application/json',
    format: 'JSON',
    data_type: 'Object',
    raw: '',
  })
}

function applyDetail(detail: InterfaceDetail) {
  form.method = detail.method
  form.url = detail.url
  form.folderId = detail.folder_id
  form.name = detail.name
  form.status = detail.status
  form.requestHeaders = toParamRows(detail.request_headers ?? [])
  form.requestBody = {
    format: detail.request_body?.format || 'json',
    data_type: detail.request_body?.data_type || 'Object',
    raw: detail.request_body?.raw ?? '',
    fields: fieldTreeFromApi(detail.request_body?.fields ?? []),
  }
  form.queryParams = toParamRows(detail.query_params ?? [])
  form.responseHeaders = toParamRows(detail.response_headers ?? [])
  form.responseResults = detail.response_results?.length
    ? detail.response_results.map((item) => ({
        ...item,
        fields: responseFieldTreeFromApi(item.fields ?? []),
      }))
    : [defaultResult()]
  form.responseExamples = detail.response_examples?.length
    ? detail.response_examples.map((item) => enrichResponseExample({ ...item }))
    : [defaultExample()]
}

async function loadDetail() {
  if (!interfaceId.value || !activeWorkspaceId.value || !activeProjectId.value) return
  loading.value = true
  try {
    const detail = await fetchInterfaceDetail(
      activeWorkspaceId.value,
      activeProjectId.value,
      interfaceId.value,
    )
    applyDetail(detail)
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('workspace.loadApiFailed'))
  } finally {
    loading.value = false
  }
}

watch(interfaceId, (id) => {
  if (id) loadDetail()
}, { immediate: true })

watch(
  () => [form.method, form.url, form.folderId, form.name, form.status, form.requestHeaders, form.requestBody, form.queryParams, form.responseHeaders, form.responseResults, form.responseExamples],
  () => {
    if (!readOnly.value) markActiveTabDirty()
  },
  { deep: true },
)

async function handleSave() {
  if (!interfaceId.value) return
  if (!form.name.trim()) {
    ElMessage.warning(t('workspace.nameRequired'))
    return
  }
  if (!form.folderId) {
    ElMessage.warning(t('workspace.folderRequired'))
    return
  }

  saving.value = true
  try {
    await submitUpdateInterface(interfaceId.value, {
      folder_id: form.folderId,
      name: form.name.trim(),
      method: form.method,
      url: form.url.trim(),
      status: form.status,
      request_headers: compactParamRows(form.requestHeaders),
      request_body: {
        format: form.requestBody.format,
        data_type: form.requestBody.data_type,
        raw: form.requestBody.raw ?? '',
        fields: compactFieldTree(form.requestBody.fields),
      },
      query_params: compactParamRows(form.queryParams),
      response_headers: compactParamRows(form.responseHeaders),
      response_results: form.responseResults.map((result) => ({
        ...result,
        fields: compactResponseFieldTree(result.fields),
      })),
      response_examples: form.responseExamples.map((item) => normalizeResponseExampleForSave(item)),
    })
    ElMessage.success(t('workspace.updateApiSuccess'))
    await loadDetail()
  } catch (error) {
    ElMessage.error(error instanceof Error ? error.message : t('workspace.updateApiFailed'))
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div v-loading="loading" class="interface-detail interface-editor">
    <div class="interface-detail__header">
      <nav class="interface-detail__subnav">
        <button
          type="button"
          class="interface-detail__subnav-item"
          :class="{ 'interface-detail__subnav-item--active': mode === 'doc' }"
          @click="mode = 'doc'"
        >
          {{ t('workspace.apiSubnav.doc') }}
        </button>
        <button
          type="button"
          class="interface-detail__subnav-item"
          :class="{ 'interface-detail__subnav-item--active': mode === 'edit' }"
          @click="mode = 'edit'"
        >
          {{ t('workspace.apiSubnav.edit') }}
        </button>
      </nav>

      <div v-if="!readOnly" class="interface-detail__toolbar">
        <el-button type="primary" :loading="saving" @click="handleSave">
          {{ t('common.save') }}
        </el-button>
      </div>
    </div>

    <div class="interface-detail__scroll">
      <div class="interface-detail__section">
      <div class="interface-detail__row">
        <InterfaceUrlBar
          v-model:protocol="form.protocol"
          v-model:method="form.method"
          v-model:url="form.url"
          :readonly="readOnly"
        />
      </div>

      <div class="interface-detail__row interface-meta-row">
        <el-select
          v-model="form.folderId"
          class="interface-meta-row__folder"
          :placeholder="t('workspace.selectFolder')"
          size="default"
          :disabled="readOnly"
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
          :readonly="readOnly"
        />
      </div>

      <div class="interface-detail__status-row">
        <span class="interface-detail__label">{{ t('workspace.columns.status') }}</span>
        <InterfaceStatusRadio v-model="form.status" :readonly="readOnly" />
      </div>
    </div>

    <div class="interface-detail__panels">
    <section class="interface-panel" :class="{ 'interface-panel--open': openPanels.requestParams }">
      <button
        type="button"
        class="interface-panel__summary"
        :aria-expanded="openPanels.requestParams"
        @click="openPanels.requestParams = !openPanels.requestParams"
      >
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.requestParams') }}
      </button>
      <div class="interface-panel__collapse">
        <div class="interface-panel__body">
          <InterfaceRequestParams
            v-model:request-headers="form.requestHeaders"
            v-model:query-params="form.queryParams"
            v-model:request-body="form.requestBody"
            :readonly="readOnly"
          />
        </div>
      </div>
    </section>

    <section class="interface-panel" :class="{ 'interface-panel--open': openPanels.responseHeader }">
      <button
        type="button"
        class="interface-panel__summary"
        :aria-expanded="openPanels.responseHeader"
        @click="openPanels.responseHeader = !openPanels.responseHeader"
      >
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.responseHeader') }}
      </button>
      <div class="interface-panel__collapse">
        <div class="interface-panel__body">
          <InterfaceResponseHeaderTable
            v-model="form.responseHeaders"
            :readonly="readOnly"
          />
        </div>
      </div>
    </section>

    <section class="interface-panel" :class="{ 'interface-panel--open': openPanels.responseResult }">
      <button
        type="button"
        class="interface-panel__summary"
        :aria-expanded="openPanels.responseResult"
        @click="openPanels.responseResult = !openPanels.responseResult"
      >
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.responseResult') }}
      </button>
      <div class="interface-panel__collapse">
        <div class="interface-panel__body">
          <InterfaceResponseResults
            v-model="form.responseResults"
            :readonly="readOnly"
          />
        </div>
      </div>
    </section>

    <section class="interface-panel" :class="{ 'interface-panel--open': openPanels.responseExample }">
      <button
        type="button"
        class="interface-panel__summary"
        :aria-expanded="openPanels.responseExample"
        @click="openPanels.responseExample = !openPanels.responseExample"
      >
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.responseExample') }}
      </button>
      <div class="interface-panel__collapse">
        <div class="interface-panel__body">
          <InterfaceResponseExamples
            v-model="form.responseExamples"
            :readonly="readOnly"
          />
        </div>
      </div>
    </section>
    </div>
    </div>
  </div>
</template>

<style scoped>
.interface-detail__subnav {
  display: flex;
  align-items: stretch;
  height: 40px;
  padding: 0 16px;
  border-bottom: 1px solid var(--color-border);
  background: var(--color-workspace-elevated);
}

.interface-detail__subnav-item {
  position: relative;
  padding: 0 16px;
  border: none;
  background: transparent;
  color: var(--color-workspace-tab-text);
  font-size: 14px;
  cursor: pointer;
}

.interface-detail__subnav-item--active {
  color: var(--color-workspace-tab-text-active);
}

.interface-detail__subnav-item--active::after {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 2px;
  background: var(--color-interface-accent);
}

.interface-detail__toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-bottom: 1px solid var(--color-border);
}

.interface-detail__section {
  display: flex;
  flex-direction: column;
  gap: 15px;
  padding: 16px var(--color-interface-section-padding-x);
  border-bottom: 1px solid var(--color-border);
}

.interface-detail__row {
  display: flex;
  align-items: center;
  width: 100%;
}

.interface-detail__status-row {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.interface-detail__label {
  flex-shrink: 0;
  color: var(--color-interface-field-text);
}
</style>
