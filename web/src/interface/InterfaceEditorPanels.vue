<script setup lang="ts">
import { reactive } from 'vue'
import { ArrowDown } from '@element-plus/icons-vue'
import InterfaceRequestParams from '@/interface/InterfaceRequestParams.vue'
import InterfaceResponseHeaderTable from '@/interface/InterfaceResponseHeaderTable.vue'
import InterfaceResponseExamples from '@/interface/InterfaceResponseExamples.vue'
import InterfaceResponseResults from '@/interface/InterfaceResponseResults.vue'
import { useLocale } from '@/composables/useLocale'
import type { ParamRow } from '@/utils/interface-params'
import type {
  InterfaceRequestBody,
  InterfaceResponseExample,
  InterfaceResponseResult,
} from '@/types/workspace'

defineProps<{
  requestHeaders: ParamRow[]
  queryParams: ParamRow[]
  requestBody: InterfaceRequestBody
  responseHeaders: ParamRow[]
  responseResults: InterfaceResponseResult[]
  responseExamples: InterfaceResponseExample[]
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:requestHeaders': [value: ParamRow[]]
  'update:queryParams': [value: ParamRow[]]
  'update:requestBody': [value: InterfaceRequestBody]
  'update:responseHeaders': [value: ParamRow[]]
  'update:responseResults': [value: InterfaceResponseResult[]]
  'update:responseExamples': [value: InterfaceResponseExample[]]
}>()

const { t } = useLocale()

const openPanels = reactive({
  requestParams: true,
  responseHeader: true,
  responseResult: true,
  responseExample: true,
})
</script>

<template>
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
            :request-headers="requestHeaders"
            :query-params="queryParams"
            :request-body="requestBody"
            :readonly="readonly"
            @update:request-headers="emit('update:requestHeaders', $event)"
            @update:query-params="emit('update:queryParams', $event)"
            @update:request-body="emit('update:requestBody', $event)"
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
            :model-value="responseHeaders"
            :readonly="readonly"
            @update:model-value="emit('update:responseHeaders', $event)"
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
            :model-value="responseResults"
            :readonly="readonly"
            @update:model-value="emit('update:responseResults', $event)"
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
            :model-value="responseExamples"
            :readonly="readonly"
            @update:model-value="emit('update:responseExamples', $event)"
          />
        </div>
      </div>
    </section>
  </div>
</template>
