<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowDown, Folder, Plus, Search } from '@element-plus/icons-vue'
import { useLocale } from '@/composables/useLocale'
import { useWorkspaceContext } from '@/composables/useWorkspaceContext'
import type { ApiTreeNode, HttpMethod } from '@/types/workspace'

const { t } = useLocale()
const {
  contextMode,
  projects,
  apiTree,
  selectedApiId,
  loadingProjects,
  selectProject,
  selectApi,
  toggleNode,
  isNodeExpanded,
  openCreateProject,
} = useWorkspaceContext()

const searchQuery = ref('')

const methodColors: Record<HttpMethod, string> = {
  GET: 'get',
  POST: 'post',
  PUT: 'put',
  DELETE: 'delete',
  PATCH: 'patch',
  HEAD: 'head',
  OPTIONS: 'options',
}

const filteredProjects = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return projects.value
  return projects.value.filter((item) => item.name.toLowerCase().includes(q))
})

function filterTree(nodes: ApiTreeNode[], query: string): ApiTreeNode[] {
  if (!query) return nodes
  const q = query.toLowerCase()

  return nodes.reduce<ApiTreeNode[]>((acc, node) => {
    if (node.type === 'api') {
      if (node.name.toLowerCase().includes(q)) acc.push(node)
      return acc
    }

    const children = filterTree(node.children ?? [], query)
    if (node.name.toLowerCase().includes(q) || children.length > 0) {
      acc.push({ ...node, children })
    }
    return acc
  }, [])
}

const filteredApiTree = computed(() => filterTree(apiTree.value, searchQuery.value.trim()))

function handleProjectClick(id: number) {
  selectProject(id)
}

function renderNodes(nodes: ApiTreeNode[], depth = 0): { node: ApiTreeNode; depth: number }[] {
  const result: { node: ApiTreeNode; depth: number }[] = []

  for (const node of nodes) {
    result.push({ node, depth })
    if (node.type === 'folder' && isNodeExpanded(node.id) && node.children?.length) {
      result.push(...renderNodes(node.children, depth + 1))
    }
  }

  return result
}

const flatApiNodes = computed(() => renderNodes(filteredApiTree.value))

function handleCreateProject() {
  openCreateProject()
}
</script>

<template>
  <aside class="context-panel">
    <div class="context-panel__header">
      <div class="context-panel__search">
        <el-icon :size="16"><Search /></el-icon>
        <input
          v-model="searchQuery"
          type="text"
          class="context-panel__search-input"
          :placeholder="
            contextMode === 'workspace'
              ? t('workspace.searchProject')
              : t('workspace.searchApi')
          "
        />
      </div>
      <button
        v-if="contextMode === 'workspace'"
        type="button"
        class="context-panel__create-btn"
        @click="handleCreateProject"
      >
        <el-icon :size="16"><Plus /></el-icon>
        <span>{{ t('workspace.createProject') }}</span>
      </button>
      <button
        v-else
        type="button"
        class="context-panel__add-btn"
      >
        <el-icon :size="16"><Plus /></el-icon>
        <span>{{ t('workspace.addApi') }}</span>
        <el-icon :size="12"><ArrowDown /></el-icon>
      </button>
    </div>

    <!-- Project list (workspace mode) -->
    <div v-if="contextMode === 'workspace'" class="context-panel__body" v-loading="loadingProjects">
      <div class="context-panel__section-title">{{ t('workspace.projectList') }}</div>
      <ul class="context-panel__list">
        <li
          v-for="item in filteredProjects"
          :key="item.id"
          class="context-panel__project-item"
          @click="handleProjectClick(item.id)"
        >
          <el-icon :size="16" class="context-panel__project-icon"><Folder /></el-icon>
          <div class="context-panel__project-info">
            <span class="context-panel__project-name">{{ item.name }}</span>
            <span v-if="item.creator_name" class="context-panel__project-desc">
              {{ item.creator_name }}
            </span>
          </div>
        </li>
        <li v-if="filteredProjects.length === 0" class="context-panel__empty">
          {{ t('workspace.noProjects') }}
        </li>
      </ul>
    </div>

    <!-- API tree (project mode) -->
    <div v-else class="context-panel__body">
      <ul class="context-panel__tree">
        <li
          v-for="{ node, depth } in flatApiNodes"
          :key="node.id"
          class="context-panel__tree-item"
          :class="{
            'context-panel__tree-item--active': node.type === 'api' && selectedApiId === node.id,
            'context-panel__tree-item--folder': node.type === 'folder',
          }"
          :style="{ paddingLeft: `${12 + depth * 16}px` }"
          @click="
            node.type === 'folder'
              ? toggleNode(node.id)
              : selectApi(node.id)
          "
        >
          <button
            v-if="node.type === 'folder'"
            type="button"
            class="context-panel__expand"
            @click.stop="toggleNode(node.id)"
          >
            <el-icon
              :size="12"
              class="context-panel__expand-icon"
              :class="{ 'context-panel__expand-icon--open': isNodeExpanded(node.id) }"
            >
              <ArrowDown />
            </el-icon>
          </button>
          <span v-else class="context-panel__expand context-panel__expand--leaf" />

          <el-icon v-if="node.type === 'folder'" :size="14" class="context-panel__folder-icon">
            <Folder />
          </el-icon>

          <span
            v-if="node.type === 'api' && node.method"
            class="context-panel__method"
            :class="`context-panel__method--${methodColors[node.method]}`"
          >
            {{ node.method }}
          </span>

          <span class="context-panel__tree-name">{{ node.name }}</span>
        </li>
        <li v-if="flatApiNodes.length === 0" class="context-panel__empty">
          {{ t('workspace.noApis') }}
        </li>
      </ul>
    </div>
  </aside>
</template>

<style scoped>
.context-panel {
  width: 280px;
  flex-shrink: 0;
  border-right: 1px solid var(--color-border);
  background: var(--color-sidebar);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.context-panel__header {
  padding: 12px 14px;
  border-bottom: 1px solid var(--color-border);
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.context-panel__search {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 12px;
  height: 38px;
  border-radius: 8px;
  background: var(--color-bg);
  border: 1px solid var(--color-border);
}

.context-panel__search-input {
  flex: 1;
  border: none;
  outline: none;
  background: transparent;
  color: var(--color-text);
  font-size: 14px;
}

.context-panel__search-input::placeholder {
  color: var(--color-text-secondary);
}

.context-panel__create-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 38px;
  padding: 0 14px;
  border: 1px dashed var(--color-border);
  border-radius: 8px;
  background: transparent;
  color: var(--color-primary-light);
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s ease, border-color 0.15s ease;
}

.context-panel__create-btn:hover {
  background: var(--color-active);
  border-color: var(--color-primary);
}

.context-panel__add-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  height: 38px;
  padding: 0 14px;
  border: none;
  border-radius: 8px;
  background: var(--color-primary);
  color: #fff;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.context-panel__add-btn:hover {
  background: var(--color-primary-light);
}

.context-panel__body {
  flex: 1;
  overflow-y: auto;
  padding: 8px 0;
}

.context-panel__section-title {
  padding: 6px 16px 10px;
  font-size: 12px;
  font-weight: 600;
  color: var(--color-text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.context-panel__list,
.context-panel__tree {
  list-style: none;
  margin: 0;
  padding: 0;
}

.context-panel__project-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 10px 16px;
  cursor: pointer;
  transition: background-color 0.15s ease;
}

.context-panel__project-item:hover {
  background: var(--color-hover);
}

.context-panel__project-icon {
  color: var(--color-primary-light);
  margin-top: 2px;
  flex-shrink: 0;
}

.context-panel__project-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.context-panel__project-name {
  font-size: 15px;
  font-weight: 500;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.context-panel__project-desc {
  font-size: 13px;
  color: var(--color-text-secondary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.context-panel__tree-item {
  display: flex;
  align-items: center;
  gap: 6px;
  height: 36px;
  padding-right: 14px;
  cursor: pointer;
  transition: background-color 0.15s ease;
  user-select: none;
}

.context-panel__tree-item:hover {
  background: var(--color-hover);
}

.context-panel__tree-item--active {
  background: var(--color-active);
}

.context-panel__tree-item--active .context-panel__tree-name {
  color: var(--color-primary-light);
  font-weight: 500;
}

.context-panel__expand {
  width: 16px;
  height: 16px;
  border: none;
  padding: 0;
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.context-panel__expand--leaf {
  cursor: default;
}

.context-panel__expand-icon {
  transition: transform 0.15s ease;
  transform: rotate(-90deg);
}

.context-panel__expand-icon--open {
  transform: rotate(0deg);
}

.context-panel__folder-icon {
  color: var(--color-text-secondary);
  flex-shrink: 0;
}

.context-panel__method {
  flex-shrink: 0;
  padding: 2px 6px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 700;
  line-height: 1.4;
  font-family: 'Consolas', 'Monaco', monospace;
}

.context-panel__method--get {
  color: #61affe;
  background: rgba(97, 175, 254, 0.12);
}

.context-panel__method--post {
  color: #49cc90;
  background: rgba(73, 204, 144, 0.12);
}

.context-panel__method--put {
  color: #fca130;
  background: rgba(252, 161, 48, 0.12);
}

.context-panel__method--delete {
  color: #f93e3e;
  background: rgba(249, 62, 62, 0.12);
}

.context-panel__method--patch {
  color: #50e3c2;
  background: rgba(80, 227, 194, 0.12);
}

.context-panel__method--head,
.context-panel__method--options {
  color: var(--color-text-secondary);
  background: var(--color-hover);
}

.context-panel__tree-name {
  font-size: 14px;
  color: var(--color-text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.context-panel__empty {
  padding: 28px 16px;
  text-align: center;
  font-size: 14px;
  color: var(--color-text-secondary);
}
</style>
