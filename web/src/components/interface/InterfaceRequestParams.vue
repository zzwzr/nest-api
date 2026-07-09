<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useLocale } from '@/composables/useLocale'
import InterfaceParamTable from '@/components/interface/InterfaceParamTable.vue'
import { compactParamRows, emptyParamRow, ensureTrailingEmptyRow, type ParamRow } from '@/utils/interface-params'
import type { InterfaceRequestBody } from '@/types/workspace'

type RequestTab = 'header' | 'body' | 'query'

const props = defineProps<{
  requestHeaders: ParamRow[]
  queryParams: ParamRow[]
  requestBody: InterfaceRequestBody
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:requestHeaders': [value: ParamRow[]]
  'update:queryParams': [value: ParamRow[]]
  'update:requestBody': [value: InterfaceRequestBody]
}>()

const { t } = useLocale()
const activeTab = ref<RequestTab>('body')

const bodyFormatOptions = [
  { value: 'form-data', label: 'Form-data' },
  { value: 'json', label: 'JSON' },
  { value: 'xml', label: 'XML' },
  { value: 'raw', label: 'Raw' },
  { value: 'binary', label: 'Binary' },
]

const tabs = computed(() => [
  { key: 'header' as const, label: t('workspace.interfaceForm.reqHeader') },
  { key: 'body' as const, label: t('workspace.interfaceForm.reqBody') },
  { key: 'query' as const, label: t('workspace.interfaceForm.queryParams') },
])

const bodyFields = computed({
  get: () => props.requestBody.fields,
  set: (fields: ParamRow[]) => {
    emit('update:requestBody', { ...props.requestBody, fields })
  },
})

function initRows(rows: ParamRow[]) {
  return ensureTrailingEmptyRow(rows.length ? rows : [emptyParamRow()])
}

watch(
  () => props.requestHeaders,
  (rows) => {
    if (!rows.length) emit('update:requestHeaders', initRows([]))
  },
  { immediate: true },
)

watch(
  () => props.queryParams,
  (rows) => {
    if (!rows.length) emit('update:queryParams', initRows([]))
  },
  { immediate: true },
)

watch(
  () => props.requestBody.fields,
  (rows) => {
    if (!rows.length) {
      emit('update:requestBody', { ...props.requestBody, fields: initRows([]) })
    }
  },
  { immediate: true },
)

function updateBodyFormat(format: string) {
  emit('update:requestBody', { ...props.requestBody, format })
}

defineExpose({
  compactHeaders: () => compactParamRows(props.requestHeaders),
  compactQuery: () => compactParamRows(props.queryParams),
  compactBodyFields: () => compactParamRows(props.requestBody.fields),
})
</script>

<template>
  <div class="interface-request-params">
    <div class="interface-request-params__tabs">
      <button
        v-for="tab in tabs"
        :key="tab.key"
        type="button"
        class="interface-request-params__tab"
        :class="{ 'interface-request-params__tab--active': activeTab === tab.key }"
        @click="activeTab = tab.key"
      >
        {{ tab.label }}
      </button>
    </div>

    <div v-if="activeTab === 'header'" class="interface-request-params__content">
      <InterfaceParamTable
        :model-value="requestHeaders"
        :readonly="readonly"
        @update:model-value="emit('update:requestHeaders', $event)"
      />
    </div>

    <div v-else-if="activeTab === 'query'" class="interface-request-params__content">
      <InterfaceParamTable
        :model-value="queryParams"
        :readonly="readonly"
        @update:model-value="emit('update:queryParams', $event)"
      />
    </div>

    <div v-else class="interface-request-params__content">
      <div class="interface-request-params__body-toolbar">
        <el-radio-group
          :model-value="requestBody.format"
          :disabled="readonly"
          @update:model-value="updateBodyFormat"
        >
          <el-radio v-for="item in bodyFormatOptions" :key="item.value" :value="item.value">
            {{ item.label }}
          </el-radio>
        </el-radio-group>
      </div>
      <InterfaceParamTable v-model="bodyFields" :readonly="readonly" />
    </div>
  </div>
</template>

