<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { ArrowDown } from '@element-plus/icons-vue'
import { fetchInterfaceDetail } from '@/api/interface'
import InterfaceRequestParams from '@/components/interface/InterfaceRequestParams.vue'
import InterfaceStatusRadio from '@/components/interface/InterfaceStatusRadio.vue'
import InterfaceUrlBar from '@/components/interface/InterfaceUrlBar.vue'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext, parseApiId } from '@/composables/useWorkspaceContext'
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
  InterfaceResponseField,
  InterfaceResponseHeader,
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
const activeResultIndex = ref(0)
const activeExampleIndex = ref(0)

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
    fields: [emptyParamRow()],
  } as InterfaceRequestBody,
  queryParams: [emptyParamRow()] as ParamRow[],
  responseHeaders: [] as InterfaceResponseHeader[],
  responseResults: [] as InterfaceResponseResult[],
  responseExamples: [] as InterfaceResponseExample[],
})

const fieldTypeOptions = ['string', 'number', 'boolean', 'object', 'array']
const formatOptions = ['JSON', 'XML', 'HTML', 'Text']
const dataTypeOptions = ['Object', 'Array', 'String', 'Number', 'Boolean']

const interfaceId = computed(() => {
  const tab = activeModuleTab.value
  if (!tab?.apiId) return null
  return parseApiId(tab.apiId)
})

const activeResult = computed(() => form.responseResults[activeResultIndex.value] ?? null)
const activeExample = computed(() => form.responseExamples[activeExampleIndex.value] ?? null)
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
    fields: [],
  }
}

function defaultExample(): InterfaceResponseExample {
  return {
    name: t('workspace.interfaceForm.successExample'),
    status_code: 200,
    content_type: 'application/json',
    raw: '{\n  \n}',
  }
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
    fields: toParamRows(detail.request_body?.fields ?? []),
  }
  form.queryParams = toParamRows(detail.query_params ?? [])
  form.responseHeaders = detail.response_headers?.length
    ? detail.response_headers.map((item) => ({ ...item }))
    : []
  form.responseResults = detail.response_results?.length
    ? detail.response_results.map((item) => ({
        ...item,
        fields: item.fields?.map((field) => ({ ...field })) ?? [],
      }))
    : [defaultResult()]
  form.responseExamples = detail.response_examples?.length
    ? detail.response_examples.map((item) => ({ ...item }))
    : [defaultExample()]
  activeResultIndex.value = 0
  activeExampleIndex.value = 0
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

function addResponseHeader() {
  form.responseHeaders.push({
    name: '',
    type: 'string',
    required: false,
    description: '',
    example: '',
  })
}

function removeResponseHeader(index: number) {
  form.responseHeaders.splice(index, 1)
}

function addResponseResult() {
  form.responseResults.push(defaultResult())
  activeResultIndex.value = form.responseResults.length - 1
}

function removeResponseResult(index: number) {
  if (form.responseResults.length <= 1) return
  form.responseResults.splice(index, 1)
  if (activeResultIndex.value >= form.responseResults.length) {
    activeResultIndex.value = form.responseResults.length - 1
  }
}

function addResponseField(fields: InterfaceResponseField[]) {
  fields.push({
    parent_id: 0,
    name: '',
    type: 'string',
    required: false,
    description: '',
    mock: '',
    example: '',
  })
}

function removeResponseField(fields: InterfaceResponseField[], index: number) {
  fields.splice(index, 1)
}

function addResponseExample() {
  form.responseExamples.push(defaultExample())
  activeExampleIndex.value = form.responseExamples.length - 1
}

function removeResponseExample(index: number) {
  if (form.responseExamples.length <= 1) return
  form.responseExamples.splice(index, 1)
  if (activeExampleIndex.value >= form.responseExamples.length) {
    activeExampleIndex.value = form.responseExamples.length - 1
  }
}

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
        fields: compactParamRows(form.requestBody.fields),
      },
      query_params: compactParamRows(form.queryParams),
      response_headers: form.responseHeaders,
      response_results: form.responseResults,
      response_examples: form.responseExamples,
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
    <details class="interface-panel" open>
      <summary class="interface-panel__summary">
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.requestParams') }}
      </summary>
      <div class="interface-panel__body">
        <InterfaceRequestParams
          v-model:request-headers="form.requestHeaders"
          v-model:query-params="form.queryParams"
          v-model:request-body="form.requestBody"
          :readonly="readOnly"
        />
      </div>
    </details>

    <details class="interface-panel" open>
      <summary class="interface-panel__summary">
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.responseHeader') }}
      </summary>
      <div class="interface-panel__body">
      <div class="interface-detail__table-wrap">
        <table class="interface-param-table">
          <thead>
            <tr>
              <th>{{ t('workspace.interfaceForm.paramName') }}</th>
              <th>{{ t('workspace.interfaceForm.paramType') }}</th>
              <th>{{ t('workspace.interfaceForm.required') }}</th>
              <th>{{ t('workspace.interfaceForm.description') }}</th>
              <th>{{ t('workspace.interfaceForm.example') }}</th>
              <th v-if="!readOnly" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, index) in form.responseHeaders" :key="index">
              <td><el-input v-model="row.name" :readonly="readOnly" /></td>
              <td>
                <el-select v-model="row.type" :disabled="readOnly" popper-class="app-action-dropdown">
                  <el-option v-for="type in fieldTypeOptions" :key="type" :label="type" :value="type" />
                </el-select>
              </td>
              <td class="interface-detail__center"><el-checkbox v-model="row.required" :disabled="readOnly" /></td>
              <td><el-input v-model="row.description" :readonly="readOnly" /></td>
              <td><el-input v-model="row.example" :readonly="readOnly" /></td>
              <td v-if="!readOnly">
                <button type="button" class="interface-detail__row-btn" @click="removeResponseHeader(index)">
                  {{ t('common.delete') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <button v-if="!readOnly" type="button" class="interface-detail__add-btn" @click="addResponseHeader">
          + {{ t('workspace.interfaceForm.addRow') }}
        </button>
      </div>
      </div>
    </details>

    <details class="interface-panel" open>
      <summary class="interface-panel__summary">
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.responseResult') }}
      </summary>
      <div class="interface-panel__body">
      <div class="interface-detail__result-tabs">
        <button
          v-for="(result, index) in form.responseResults"
          :key="index"
          type="button"
          class="interface-detail__result-tab"
          :class="{ 'interface-detail__result-tab--active': activeResultIndex === index }"
          @click="activeResultIndex = index"
        >
          {{ result.name || t('workspace.interfaceForm.responseResult') }} ({{ result.status_code }})
        </button>
        <button v-if="!readOnly" type="button" class="interface-detail__result-tab-add" @click="addResponseResult">+</button>
      </div>

      <div v-if="activeResult" class="interface-detail__result-meta">
        <el-input v-model="activeResult.name" :readonly="readOnly" class="interface-detail__result-name" />
        <el-input-number v-model="activeResult.status_code" :min="100" :max="599" :disabled="readOnly" controls-position="right" />
        <el-select v-model="activeResult.format" :disabled="readOnly" popper-class="app-action-dropdown">
          <el-option v-for="item in formatOptions" :key="item" :label="item" :value="item" />
        </el-select>
        <el-select v-model="activeResult.data_type" :disabled="readOnly" popper-class="app-action-dropdown">
          <el-option v-for="item in dataTypeOptions" :key="item" :label="item" :value="item" />
        </el-select>
        <button v-if="!readOnly && form.responseResults.length > 1" type="button" class="interface-detail__row-btn" @click="removeResponseResult(activeResultIndex)">
          {{ t('common.delete') }}
        </button>
      </div>

      <div v-if="activeResult" class="interface-detail__table-wrap">
        <table class="interface-param-table">
          <thead>
            <tr>
              <th>{{ t('workspace.interfaceForm.paramName') }}</th>
              <th>{{ t('workspace.interfaceForm.paramType') }}</th>
              <th>{{ t('workspace.interfaceForm.required') }}</th>
              <th>{{ t('workspace.interfaceForm.description') }}</th>
              <th>Mock</th>
              <th>{{ t('workspace.interfaceForm.example') }}</th>
              <th v-if="!readOnly" />
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, index) in activeResult.fields" :key="index">
              <td><el-input v-model="row.name" :readonly="readOnly" /></td>
              <td>
                <el-select v-model="row.type" :disabled="readOnly" popper-class="app-action-dropdown">
                  <el-option v-for="type in fieldTypeOptions" :key="type" :label="type" :value="type" />
                </el-select>
              </td>
              <td class="interface-detail__center"><el-checkbox v-model="row.required" :disabled="readOnly" /></td>
              <td><el-input v-model="row.description" :readonly="readOnly" /></td>
              <td><el-input v-model="row.mock" :readonly="readOnly" /></td>
              <td><el-input v-model="row.example" :readonly="readOnly" /></td>
              <td v-if="!readOnly">
                <button type="button" class="interface-detail__row-btn" @click="removeResponseField(activeResult.fields, index)">
                  {{ t('common.delete') }}
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <button v-if="!readOnly && activeResult" type="button" class="interface-detail__add-btn" @click="addResponseField(activeResult.fields)">
          + {{ t('workspace.interfaceForm.addRow') }}
        </button>
      </div>
      </div>
    </details>

    <details class="interface-panel" open>
      <summary class="interface-panel__summary">
        <span class="interface-panel__arrow"><el-icon :size="12"><ArrowDown /></el-icon></span>
        {{ t('workspace.interfaceForm.responseExample') }}
      </summary>
      <div class="interface-panel__body">
      <div class="interface-detail__result-tabs">
        <button
          v-for="(example, index) in form.responseExamples"
          :key="index"
          type="button"
          class="interface-detail__result-tab"
          :class="{ 'interface-detail__result-tab--active': activeExampleIndex === index }"
          @click="activeExampleIndex = index"
        >
          {{ example.name }}
        </button>
        <button v-if="!readOnly" type="button" class="interface-detail__result-tab-add" @click="addResponseExample">+</button>
      </div>

      <div v-if="activeExample" class="interface-detail__example-meta">
        <el-input v-model="activeExample.name" :readonly="readOnly" class="interface-detail__result-name" />
        <el-input-number v-model="activeExample.status_code" :min="100" :max="599" :disabled="readOnly" controls-position="right" />
        <el-input v-model="activeExample.content_type" :readonly="readOnly" class="interface-detail__content-type" />
        <button v-if="!readOnly && form.responseExamples.length > 1" type="button" class="interface-detail__row-btn" @click="removeResponseExample(activeExampleIndex)">
          {{ t('common.delete') }}
        </button>
      </div>

      <el-input v-if="activeExample" v-model="activeExample.raw" type="textarea" :rows="12" class="interface-detail__raw" :readonly="readOnly" />
      </div>
    </details>
    </div>
  </div>
</template>

<style scoped>
.interface-detail {
  padding: 0 0 32px;
}

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
  padding: 12px 16px 0;
}

.interface-detail__section {
  display: flex;
  flex-direction: column;
  gap: 10px;
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

.interface-detail__table-wrap {
  padding: 12px 12px 0;
}

.interface-detail__center {
  text-align: center;
}

.interface-detail__add-btn,
.interface-detail__row-btn {
  margin-top: 8px;
  border: none;
  background: transparent;
  color: var(--color-interface-accent);
  cursor: pointer;
  font-size: 14px;
}

.interface-detail__row-btn {
  margin-top: 0;
  color: var(--color-danger);
}

.interface-detail__result-tabs {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 12px 12px 12px;
  flex-wrap: wrap;
}

.interface-detail__result-tab,
.interface-detail__result-tab-add {
  border: none;
  background: transparent;
  color: var(--color-workspace-tab-text);
  padding: 6px 10px;
  cursor: pointer;
  font-size: 14px;
}

.interface-detail__result-meta,
.interface-detail__example-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 12px 12px;
  flex-wrap: wrap;
}

.interface-detail__result-name {
  width: 160px;
}

.interface-detail__content-type {
  width: 220px;
}

.interface-detail__raw {
  margin: 0 12px;
  width: calc(100% - 24px);
}
</style>
