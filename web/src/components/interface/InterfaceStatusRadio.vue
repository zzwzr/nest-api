<script setup lang="ts">
import { computed } from 'vue'
import { INTERFACE_STATUS_OPTIONS } from '@/constants/interface-status'
import { useLocale } from '@/composables/useLocale'
import type { InterfaceStatus } from '@/types/workspace'

defineProps<{
  modelValue: InterfaceStatus
  readonly?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: InterfaceStatus]
}>()

const { t } = useLocale()

const statusOptions = computed(() =>
  INTERFACE_STATUS_OPTIONS.map((item) => ({
    value: item.value,
    label: t(item.labelKey),
    color: item.color,
  })),
)
</script>

<template>
  <el-radio-group
    :model-value="modelValue"
    class="interface-status-radio"
    :class="{ 'interface-status-radio--readonly': readonly }"
    :disabled="readonly"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <el-radio
      v-for="item in statusOptions"
      :key="item.value"
      :value="item.value"
      class="interface-status-radio__item"
    >
      <span class="interface-status-radio__label" :style="{ color: item.color }">
        {{ item.label }}
      </span>
    </el-radio>
  </el-radio-group>
</template>

<style scoped>
.interface-status-radio {
  display: flex;
  flex-wrap: wrap;
  gap: 2px 14px;
}

.interface-status-radio__item {
  margin-right: 0;
}

.interface-status-radio__label {
  font-size: 14px;
}
</style>
