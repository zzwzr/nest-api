<script setup lang="ts">
import { watch } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import { useLocale } from '@/composables/useLocale'
import {
  compactParamRows,
  emptyParamRow,
  ensureTrailingEmptyRow,
  hasParamContent,
  type ParamRow,
} from '@/utils/interface-params'

const props = defineProps<{
  modelValue: ParamRow[]
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: ParamRow[]]
}>()

const { t } = useLocale()

const fieldTypeOptions = ['string', 'number', 'boolean', 'object', 'array']

function syncRows(rows: ParamRow[]) {
  const normalized = props.readonly ? rows : ensureTrailingEmptyRow(rows)
  emit('update:modelValue', normalized)
}

function updateRow(index: number, patch: Partial<ParamRow>) {
  const rows = props.modelValue.map((row, i) => (i === index ? { ...row, ...patch } : row))
  syncRows(rows)
}

function removeRow(index: number) {
  const rows = props.modelValue.filter((_, i) => i !== index)
  syncRows(rows.length ? rows : [emptyParamRow()])
}

watch(
  () => props.modelValue,
  (rows) => {
    if (props.readonly || !rows.length) return
    const last = rows[rows.length - 1]
    if (hasParamContent(last)) {
      syncRows(rows)
    }
  },
  { deep: true },
)

defineExpose({ compact: () => compactParamRows(props.modelValue) })
</script>

<template>
  <table class="interface-param-table">
    <thead>
      <tr>
        <th>{{ t('workspace.interfaceForm.paramName') }}</th>
        <th>{{ t('workspace.interfaceForm.paramType') }}</th>
        <th class="interface-param-table__center">{{ t('workspace.interfaceForm.required') }}</th>
        <th>{{ t('workspace.interfaceForm.description') }}</th>
        <th>{{ t('workspace.interfaceForm.example') }}</th>
        <th v-if="!readonly" class="interface-param-table__actions" />
      </tr>
    </thead>
    <tbody>
      <tr v-for="(row, index) in modelValue" :key="index">
        <td>
          <el-input
            :model-value="row.name"
            :placeholder="t('workspace.interfaceForm.paramName')"
            :readonly="readonly"
            @update:model-value="updateRow(index, { name: $event })"
          />
        </td>
        <td>
          <el-select
            :model-value="row.type"
            :disabled="readonly"
            popper-class="app-action-dropdown"
            @update:model-value="updateRow(index, { type: $event })"
          >
            <el-option v-for="type in fieldTypeOptions" :key="type" :label="type" :value="type" />
          </el-select>
        </td>
        <td class="interface-param-table__center">
          <el-checkbox
            :model-value="row.required"
            :disabled="readonly"
            @update:model-value="updateRow(index, { required: $event })"
          />
        </td>
        <td>
          <el-input
            :model-value="row.description"
            :placeholder="t('workspace.interfaceForm.paramDescPlaceholder')"
            :readonly="readonly"
            @update:model-value="updateRow(index, { description: $event })"
          />
        </td>
        <td>
          <el-input
            :model-value="row.example"
            :placeholder="t('workspace.interfaceForm.example')"
            :readonly="readonly"
            @update:model-value="updateRow(index, { example: $event })"
          />
        </td>
        <td v-if="!readonly" class="interface-param-table__actions">
          <button
            v-if="hasParamContent(row)"
            type="button"
            class="interface-param-table__delete"
            @click="removeRow(index)"
          >
            <el-icon :size="14"><Delete /></el-icon>
          </button>
        </td>
      </tr>
    </tbody>
  </table>
</template>
