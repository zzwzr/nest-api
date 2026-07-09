import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import {
  createFolder,
  deleteFolder,
  fetchFolderTree,
  updateFolder,
} from '@/api/folder'
import {
  createInterface,
  deleteInterface,
  fetchInterfaces,
  reorderInterfaces,
  updateInterface,
} from '@/api/interface'
import { fetchEnvironments } from '@/api/environment'
import { fetchEnvironmentVariables } from '@/api/envvariable'
import { createProject as apiCreateProject, fetchProjects } from '@/api/project'
import { createWorkspace as apiCreateWorkspace, fetchWorkspaces } from '@/api/workspace'
import { readWorkspaceLayout, writeWorkspaceLayout } from '@/utils/workspace-layout-storage'
import type { StoredWorkspaceTab } from '@/utils/workspace-layout-storage'
import type {
  InterfaceItem,
  InterfaceDetail,
  ApiTreeNode,
  AppModule,
  EnvironmentItem,
  EnvironmentVariableItem,
  HttpMethod,
  InterfaceStatus,
  ProjectItem,
  WorkspaceItem,
  WorkspaceTab,
} from '@/types/workspace'

const storedLayout = readWorkspaceLayout()

const workspaces = ref<WorkspaceItem[]>([])
const projects = ref<ProjectItem[]>([])
const apiTree = ref<ApiTreeNode[]>([])
const interfaces = ref<InterfaceItem[]>([])
const environments = ref<EnvironmentItem[]>([])
const environmentVariables = ref<EnvironmentVariableItem[]>([])
const loadingWorkspaces = ref(false)
const loadingProjects = ref(false)
const loadingApiTree = ref(false)
const loadingInterfaces = ref(false)
const loadingEnvironments = ref(false)
const loadingEnvironmentVariables = ref(false)

const activeModule = ref<AppModule>(
  (storedLayout.activeModule as AppModule) || 'api',
)
const activeWorkspaceId = ref<number | null>(storedLayout.workspaceId)
const activeProjectId = ref<number | null>(storedLayout.projectId)
const expandedNodeIds = ref<Set<string>>(new Set())
const selectedFolderId = ref<string | null>(null)
const selectedApiId = ref<string | null>(null)
const activeEnvironmentId = ref<number | null>(null)

const createWorkspaceVisible = ref(false)
const createProjectVisible = ref(false)
const creatingApi = ref(false)

const CREATE_API_TAB_ID = 'create-api'

const workspaceTabs = ref<WorkspaceTab[]>([])
const activeTabId = ref<string | null>(storedLayout.activeTabId ?? null)
let navigationHydrated = false
let bootstrapComplete = false

export function parseFolderId(nodeId: string): number | null {
  const match = nodeId.match(/^folder-(\d+)$/)
  return match ? Number(match[1]) : null
}

export function parseApiId(nodeId: string): number | null {
  const match = nodeId.match(/^api-(\d+)$/)
  return match ? Number(match[1]) : null
}

function findNodeById(nodes: ApiTreeNode[], id: string): ApiTreeNode | null {
  for (const node of nodes) {
    if (node.id === id) return node
    if (node.children?.length) {
      const found = findNodeById(node.children, id)
      if (found) return found
    }
  }
  return null
}

function getTabModule(tab: WorkspaceTab): AppModule {
  if (tab.module) return tab.module
  switch (tab.kind) {
    case 'folder':
    case 'api':
    case 'create-api':
      return 'api'
    case 'workspace-list':
    case 'project-list':
    case 'member-list':
      return 'project'
    case 'env-list':
    case 'env-variables':
      return 'environment'
    case 'quick-test':
      return 'quick-test'
    default:
      return 'api'
  }
}

function findPathToNode(
  nodes: ApiTreeNode[],
  targetId: string,
  trail: string[] = [],
): string[] | null {
  for (const node of nodes) {
    const nextTrail = [...trail, node.id]
    if (node.id === targetId) return nextTrail
    if (node.children?.length) {
      const found = findPathToNode(node.children, targetId, nextTrail)
      if (found) return found
    }
  }
  return null
}

function expandToNode(nodes: ApiTreeNode[], nodeId: string, expanded: Set<string>) {
  const path = findPathToNode(nodes, nodeId)
  if (!path) return
  for (const id of path) {
    if (id.startsWith('folder-')) {
      expanded.add(id)
    }
  }
}

export function useWorkspaceContext() {
  const activeWorkspace = computed(
    () => workspaces.value.find((item) => item.id === activeWorkspaceId.value) ?? null,
  )

  const workspaceProjects = computed(() =>
    projects.value.filter((item) => item.workspace_id === activeWorkspaceId.value),
  )

  const activeProject = computed(
    () => workspaceProjects.value.find((item) => item.id === activeProjectId.value) ?? null,
  )

  const contextMode = computed<'workspace' | 'project'>(() =>
    activeProjectId.value ? 'project' : 'workspace',
  )

  const selectedFolder = computed(
    () =>
      (selectedFolderId.value
        ? findNodeById(apiTree.value, selectedFolderId.value)
        : null),
  )

  const folderInterfaces = computed(() => {
    if (!selectedFolderId.value) return []
    const folderId = parseFolderId(selectedFolderId.value)
    if (!folderId) return []
    return interfaces.value.filter((item) => item.folder_id === folderId)
  })

  const activeEnvironment = computed(
    () => environments.value.find((item) => item.id === activeEnvironmentId.value) ?? null,
  )

  const activeTab = computed(
    () => workspaceTabs.value.find((item) => item.id === activeTabId.value) ?? null,
  )

  const activeModuleTab = computed(() => {
    const tab = workspaceTabs.value.find((item) => item.id === activeTabId.value)
    if (!tab || getTabModule(tab) !== activeModule.value) return null
    return tab
  })

  const moduleTabs = computed(() => {
    const tabs = workspaceTabs.value.filter((item) => getTabModule(item) === activeModule.value)
    if (activeModule.value !== 'environment') return tabs

    if (activeModuleTab.value?.kind === 'env-variables') {
      return tabs.filter((item) => item.kind === 'env-variables')
    }
    return tabs.filter((item) => item.kind === 'env-list')
  })

  function resetWorkspaceTabs() {
    workspaceTabs.value = []
    activeTabId.value = null
    selectedFolderId.value = null
    selectedApiId.value = null
    creatingApi.value = false
  }

  function clearTabSelection() {
    activeTabId.value = null
    selectedFolderId.value = null
    selectedApiId.value = null
    creatingApi.value = false
  }

  function applyTabSelection(tab: WorkspaceTab) {
    creatingApi.value = false
    selectedApiId.value = null
    selectedFolderId.value = null

    switch (tab.kind) {
      case 'folder':
        selectedFolderId.value = tab.folderId ?? tab.id
        break
      case 'api':
        selectedApiId.value = tab.apiId ?? tab.id
        break
      case 'create-api':
        creatingApi.value = true
        selectedFolderId.value = tab.folderId ?? null
        break
      case 'env-variables':
        if (tab.environmentId) {
          activeEnvironmentId.value = tab.environmentId
          loadEnvironmentVariables()
        }
        break
      default:
        break
    }
  }

  function activateTab(tabId: string) {
    const tab = workspaceTabs.value.find((item) => item.id === tabId)
    if (!tab) return
    activeModule.value = getTabModule(tab)
    activeTabId.value = tabId
    applyTabSelection(tab)
    expandSelectionPath(tab)
    persistLayoutState()
  }

  function upsertTab(tab: WorkspaceTab) {
    const normalizedTab = { ...tab, module: tab.module ?? getTabModule(tab) }
    const index = workspaceTabs.value.findIndex((item) => item.id === normalizedTab.id)
    if (index === -1) {
      workspaceTabs.value.push(normalizedTab)
      return
    }
    workspaceTabs.value[index] = { ...workspaceTabs.value[index], ...normalizedTab }
  }

  function openFolderTab(id: string, name: string) {
    upsertTab({
      id,
      kind: 'folder',
      module: 'api',
      label: name,
      folderId: id,
      closable: true,
    })
    activateTab(id)
  }

  function openApiTab(id: string, name: string, method?: HttpMethod) {
    upsertTab({
      id,
      kind: 'api',
      module: 'api',
      label: name,
      apiId: id,
      method,
      closable: true,
    })
    activateTab(id)
  }

  function openCreateApiTab(folderId?: string | null) {
    const resolvedFolderId = folderId ?? selectedFolderId.value
    const folderNode = resolvedFolderId
      ? findNodeById(apiTree.value, resolvedFolderId)
      : null

    upsertTab({
      id: CREATE_API_TAB_ID,
      kind: 'create-api',
      module: 'api',
      label: '',
      folderId: resolvedFolderId ?? undefined,
      closable: true,
    })
    activateTab(CREATE_API_TAB_ID)

    if (folderNode) {
      selectedFolderId.value = folderNode.id
    }
  }

  function openSpaceSectionTab(section: 'workspaces' | 'projects' | 'members') {
    const tabMap = {
      workspaces: { id: 'workspace-list', kind: 'workspace-list' as const },
      projects: { id: 'project-list', kind: 'project-list' as const },
      members: { id: 'member-list', kind: 'member-list' as const },
    }
    const tab = tabMap[section]
    upsertTab({
      id: tab.id,
      kind: tab.kind,
      module: 'project',
      label: tab.id,
      closable: false,
    })
    activateTab(tab.id)
  }

  function openProjectSectionTab(section: 'projects' | 'members') {
    openSpaceSectionTab(section)
  }

  function openEnvSectionTab(section: 'list' | 'variables') {
    if (section === 'list') {
      upsertTab({
        id: 'env-list',
        kind: 'env-list',
        module: 'environment',
        label: 'env-list',
        closable: false,
      })
      activateTab('env-list')
      return
    }

    if (!environments.value.length) {
      upsertTab({
        id: 'env-variables-empty',
        kind: 'env-variables',
        module: 'environment',
        label: 'env-variables-empty',
        closable: false,
      })
      activateTab('env-variables-empty')
      return
    }

    for (const env of environments.value) {
      upsertTab({
        id: `env-var-${env.id}`,
        kind: 'env-variables',
        module: 'environment',
        label: env.name,
        environmentId: env.id,
        closable: false,
      })
    }

    const targetEnvId = activeEnvironmentId.value ?? environments.value[0].id
    activateTab(`env-var-${targetEnvId}`)
  }

  function openEnvVariableTab(environmentId: number, name: string) {
    const id = `env-var-${environmentId}`
    upsertTab({
      id,
      kind: 'env-variables',
      module: 'environment',
      label: name,
      environmentId,
      closable: false,
    })
    activateTab(id)
  }

  function openQuickTestTab() {
    upsertTab({
      id: 'quick-test',
      kind: 'quick-test',
      module: 'quick-test',
      label: 'quick-test',
      closable: false,
    })
    activateTab('quick-test')
  }

  function ensureModuleDefaultTab(module: AppModule) {
    const tabs = workspaceTabs.value.filter((item) => getTabModule(item) === module)
    if (tabs.length > 0) {
      const activeInModule = tabs.find((item) => item.id === activeTabId.value)
      if (!activeInModule) {
        activateTab(tabs[tabs.length - 1].id)
      }
      return
    }

    switch (module) {
      case 'project':
        openSpaceSectionTab('workspaces')
        break
      case 'environment':
        openEnvSectionTab('list')
        break
      case 'quick-test':
        openQuickTestTab()
        break
      default:
        break
    }
  }

  function closeTab(tabId: string) {
    const tab = workspaceTabs.value.find((item) => item.id === tabId)
    if (!tab || tab.closable === false) return

    const index = workspaceTabs.value.findIndex((item) => item.id === tabId)
    if (index === -1) return

    const wasActive = activeTabId.value === tabId
    workspaceTabs.value.splice(index, 1)

    if (!wasActive) return

    const nextTab =
      workspaceTabs.value[index] ??
      workspaceTabs.value[index - 1]
    if (nextTab) {
      activateTab(nextTab.id)
    } else {
      clearTabSelection()
    }
  }

  function closeTabsLeft(tabId: string) {
    const index = workspaceTabs.value.findIndex((item) => item.id === tabId)
    if (index <= 0) return

    const removedActive = workspaceTabs.value
      .slice(0, index)
      .some((item) => item.id === activeTabId.value)

    workspaceTabs.value = workspaceTabs.value.slice(index)

    if (removedActive) {
      activateTab(tabId)
    }
  }

  function closeTabsRight(tabId: string) {
    const index = workspaceTabs.value.findIndex((item) => item.id === tabId)
    if (index === -1 || index >= workspaceTabs.value.length - 1) return

    const removedActive = workspaceTabs.value
      .slice(index + 1)
      .some((item) => item.id === activeTabId.value)

    workspaceTabs.value = workspaceTabs.value.slice(0, index + 1)

    if (removedActive) {
      activateTab(tabId)
    }
  }

  function closeOtherTabs(tabId: string) {
    if (!workspaceTabs.value.some((item) => item.id === tabId)) return

    workspaceTabs.value = workspaceTabs.value.filter(
      (item) => item.id === tabId || item.closable === false,
    )
    activateTab(tabId)
  }

  function closeAllTabs() {
    workspaceTabs.value = []
    clearTabSelection()
  }

  function markActiveTabDirty() {
    const tab = workspaceTabs.value.find((item) => item.id === activeTabId.value)
    if (tab) {
      tab.dirty = true
    }
  }

  function syncTabLabels() {
    for (const tab of workspaceTabs.value) {
      if (tab.kind === 'create-api') continue

      if (tab.kind === 'env-variables' && tab.environmentId) {
        const env = environments.value.find((item) => item.id === tab.environmentId)
        if (env) tab.label = env.name
        continue
      }

      const node = findNodeById(apiTree.value, tab.id)
      if (!node) continue

      tab.label = node.name
      if (node.type === 'api') {
        tab.method = node.method
      }
    }
  }

  async function loadEnvironmentVariables() {
    if (!activeWorkspaceId.value || !activeProjectId.value || !activeEnvironmentId.value) {
      environmentVariables.value = []
      return
    }

    loadingEnvironmentVariables.value = true
    try {
      environmentVariables.value = await fetchEnvironmentVariables(
        activeWorkspaceId.value,
        activeProjectId.value,
        activeEnvironmentId.value,
      )
    } finally {
      loadingEnvironmentVariables.value = false
    }
  }

  async function loadEnvironments() {
    if (!activeWorkspaceId.value || !activeProjectId.value) {
      environments.value = []
      activeEnvironmentId.value = null
      return
    }

    loadingEnvironments.value = true
    try {
      const list = await fetchEnvironments(activeWorkspaceId.value, activeProjectId.value)
      environments.value = list

      if (!list.length) {
        activeEnvironmentId.value = null
        return
      }

      const current = list.find((item) => item.id === activeEnvironmentId.value)
      if (current) return

      const defaultEnv = list.find((item) => item.is_default)
      activeEnvironmentId.value = defaultEnv?.id ?? list[0].id
    } finally {
      loadingEnvironments.value = false
    }
  }

  async function refreshEnvironments() {
    await loadEnvironments()
  }

  async function loadWorkspaces() {
    loadingWorkspaces.value = true
    try {
      const list = await fetchWorkspaces()
      workspaces.value = list

      if (!list.length) {
        activeWorkspaceId.value = null
        activeProjectId.value = null
        projects.value = []
        apiTree.value = []
        interfaces.value = []
        return
      }

      if (!activeWorkspaceId.value || !list.some((item) => item.id === activeWorkspaceId.value)) {
        activeWorkspaceId.value = list[0].id
      }
    } finally {
      loadingWorkspaces.value = false
    }
  }

  async function loadProjects() {
    if (!activeWorkspaceId.value) {
      projects.value = []
      return
    }

    loadingProjects.value = true
    try {
      const list = await fetchProjects(activeWorkspaceId.value)
      projects.value = list

      if (
        activeProjectId.value &&
        !list.some((item) => item.id === activeProjectId.value)
      ) {
        activeProjectId.value = null
        resetWorkspaceTabs()
        apiTree.value = []
        interfaces.value = []
      }
    } finally {
      loadingProjects.value = false
    }
  }

  async function loadApiTree() {
    if (!activeWorkspaceId.value || !activeProjectId.value) {
      apiTree.value = []
      return
    }

    loadingApiTree.value = true
    try {
      apiTree.value = await fetchFolderTree(activeWorkspaceId.value, activeProjectId.value)
    } finally {
      loadingApiTree.value = false
    }

    const active = workspaceTabs.value.find((item) => item.id === activeTabId.value)
    if (active) {
      expandSelectionPath(active)
    }
  }

  async function loadInterfaces() {
    if (!activeWorkspaceId.value || !activeProjectId.value) {
      interfaces.value = []
      return
    }

    loadingInterfaces.value = true
    try {
      interfaces.value = await fetchInterfaces(
        activeWorkspaceId.value,
        activeProjectId.value,
      )
    } finally {
      loadingInterfaces.value = false
    }
  }

  async function refreshApiData() {
    await Promise.all([loadApiTree(), loadInterfaces()])
  }

  function expandSelectionPath(tab: WorkspaceTab) {
    const next = new Set(expandedNodeIds.value)
    if (tab.kind === 'folder') {
      expandToNode(apiTree.value, tab.folderId ?? tab.id, next)
    } else if (tab.kind === 'api') {
      expandToNode(apiTree.value, tab.apiId ?? tab.id, next)
    }
    expandedNodeIds.value = next
  }

  function rebuildTabFromId(tabId: string): WorkspaceTab | null {
    if (tabId === CREATE_API_TAB_ID) {
      return {
        id: CREATE_API_TAB_ID,
        kind: 'create-api',
        module: 'api',
        label: '',
        closable: true,
      }
    }

    const node = findNodeById(apiTree.value, tabId)
    if (!node) return null

    if (node.type === 'folder') {
      return {
        id: tabId,
        kind: 'folder',
        module: 'api',
        label: node.name,
        folderId: tabId,
        closable: true,
      }
    }

    return {
      id: tabId,
      kind: 'api',
      module: 'api',
      label: node.name,
      apiId: tabId,
      method: node.method,
      closable: true,
    }
  }

  function persistLayoutState() {
    if (!navigationHydrated) return

    const tabs: StoredWorkspaceTab[] = workspaceTabs.value.map((tab) => ({
      id: tab.id,
      kind: tab.kind,
      module: getTabModule(tab),
      label: tab.label,
      method: tab.method,
      folderId: tab.folderId,
      apiId: tab.apiId,
      environmentId: tab.environmentId,
      closable: tab.closable,
    }))

    writeWorkspaceLayout({
      workspaceId: activeWorkspaceId.value,
      projectId: activeProjectId.value,
      activeModule: activeModule.value,
      activeTabId: activeTabId.value,
      selectedFolderId: selectedFolderId.value,
      selectedApiId: selectedApiId.value,
      workspaceTabs: tabs,
    })
  }

  function restoreNavigationState() {
    const layout = readWorkspaceLayout()

    if (layout.activeModule) {
      activeModule.value = layout.activeModule as AppModule
    }

    const restoredTabs: WorkspaceTab[] = []
    const storedTabs = layout.workspaceTabs ?? []

    for (const stored of storedTabs) {
      const tab = stored as WorkspaceTab
      const module = getTabModule(tab)

      if (module === 'api' && tab.kind !== 'create-api') {
        if (!findNodeById(apiTree.value, tab.id)) continue
      }

      if (module === 'environment' && tab.kind === 'env-variables' && tab.environmentId) {
        if (!environments.value.some((item) => item.id === tab.environmentId)) continue
      }

      restoredTabs.push({ ...tab, module })
    }

    const fallbackIds = [
      layout.activeTabId,
      layout.selectedApiId,
      layout.selectedFolderId,
    ].filter((id): id is string => !!id)

    for (const tabId of fallbackIds) {
      if (restoredTabs.some((item) => item.id === tabId)) continue
      const rebuilt = rebuildTabFromId(tabId)
      if (rebuilt) {
        restoredTabs.push(rebuilt)
        break
      }
    }

    if (!restoredTabs.length) {
      ensureModuleDefaultTab(activeModule.value)
      return
    }

    workspaceTabs.value = restoredTabs
    syncTabLabels()

    const targetId =
      layout.activeTabId && restoredTabs.some((item) => item.id === layout.activeTabId)
        ? layout.activeTabId
        : restoredTabs[restoredTabs.length - 1].id

    activateTab(targetId)

    const active = workspaceTabs.value.find((item) => item.id === targetId)
    if (active) {
      expandSelectionPath(active)
    }
  }

  async function bootstrap() {
    await loadWorkspaces()
    await loadProjects()
    await refreshApiData()
    await loadEnvironments()
    restoreNavigationState()
    bootstrapComplete = true
    navigationHydrated = true
    persistLayoutState()
  }

  async function refreshWorkspaces() {
    await loadWorkspaces()
    await loadProjects()
    await refreshApiData()
    await loadEnvironments()
  }

  async function refreshProjects() {
    await loadProjects()
    await refreshApiData()
    await loadEnvironments()
  }

  function switchModule(module: AppModule) {
    activeModule.value = module
    ensureModuleDefaultTab(module)
  }

  function setModule(module: AppModule) {
    switchModule(module)
  }

  function selectWorkspace(id: number) {
    if (activeWorkspaceId.value === id) return
    activeWorkspaceId.value = id
    activeProjectId.value = null
    resetWorkspaceTabs()
    expandedNodeIds.value = new Set()
    apiTree.value = []
    interfaces.value = []
    environments.value = []
    activeEnvironmentId.value = null
    persistLayoutState()
  }

  function selectProject(id: number | null) {
    if (activeProjectId.value === id) return
    activeProjectId.value = id
    resetWorkspaceTabs()
    expandedNodeIds.value = new Set()
    activeEnvironmentId.value = null
    if (!id) {
      apiTree.value = []
      interfaces.value = []
      environments.value = []
    }
    persistLayoutState()
  }

  function selectEnvironment(id: number) {
    activeEnvironmentId.value = id
  }

  function selectFolder(id: string) {
    const node = findNodeById(apiTree.value, id)
    openFolderTab(id, node?.name ?? id)
  }

  function selectApi(id: string) {
    const node = findNodeById(apiTree.value, id)
    openApiTab(id, node?.name ?? id, node?.method)
  }

  function openCreateApi(folderId?: string | null) {
    openCreateApiTab(folderId)
  }

  function closeCreateApi() {
    closeTab(CREATE_API_TAB_ID)
  }

  function toggleNode(id: string) {
    const next = new Set(expandedNodeIds.value)
    if (next.has(id)) {
      next.delete(id)
    } else {
      next.add(id)
    }
    expandedNodeIds.value = next
  }

  function isNodeExpanded(id: string) {
    return expandedNodeIds.value.has(id)
  }

  function openCreateWorkspace() {
    createWorkspaceVisible.value = true
  }

  function openCreateProject() {
    if (!activeWorkspaceId.value) {
      ElMessage.warning('请先选择工作空间')
      return
    }
    createProjectVisible.value = true
  }

  async function submitCreateWorkspace(name: string) {
    await apiCreateWorkspace(name)
    createWorkspaceVisible.value = false
    await refreshWorkspaces()
  }

  async function submitCreateProject(name: string) {
    if (!activeWorkspaceId.value) return
    await apiCreateProject(activeWorkspaceId.value, name)
    createProjectVisible.value = false
    await refreshProjects()
  }

  async function submitCreateFolder(name: string, parentId = 0) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return
    await createFolder(activeWorkspaceId.value, activeProjectId.value, name, parentId)
    await refreshApiData()
  }

  async function submitUpdateFolder(folderId: number, name: string) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return
    await updateFolder(activeWorkspaceId.value, activeProjectId.value, folderId, name)
    await refreshApiData()
  }

  async function submitDeleteFolder(folderId: number) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return
    await deleteFolder(activeWorkspaceId.value, activeProjectId.value, folderId)
    if (selectedFolderId.value === `folder-${folderId}`) {
      closeTab(`folder-${folderId}`)
    }
    await refreshApiData()
  }

  async function submitCreateInterface(
    folderId: number,
    name: string,
    method: HttpMethod,
    url = '',
    status: InterfaceStatus = 1,
  ) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return
    await createInterface(activeWorkspaceId.value, activeProjectId.value, {
      folder_id: folderId,
      name,
      method,
      url,
      status,
    })
    creatingApi.value = false
    await refreshApiData()
    syncTabLabels()
    closeTab(CREATE_API_TAB_ID)
    openFolderTab(`folder-${folderId}`, findNodeById(apiTree.value, `folder-${folderId}`)?.name ?? '')
  }

  async function submitUpdateInterface(
    interfaceId: number,
    payload: {
      folder_id?: number
      name: string
      method: HttpMethod
      url?: string
      status?: number
      request_headers?: InterfaceDetail['request_headers']
      request_body?: InterfaceDetail['request_body']
      query_params?: InterfaceDetail['query_params']
      response_headers?: InterfaceDetail['response_headers']
      response_results?: InterfaceDetail['response_results']
      response_examples?: InterfaceDetail['response_examples']
    },
  ) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return
    await updateInterface(
      activeWorkspaceId.value,
      activeProjectId.value,
      interfaceId,
      payload,
    )
    const tab = workspaceTabs.value.find((item) => item.id === `api-${interfaceId}`)
    if (tab) tab.dirty = false
    await refreshApiData()
    syncTabLabels()
  }

  async function submitDeleteInterface(interfaceId: number) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return
    await deleteInterface(activeWorkspaceId.value, activeProjectId.value, interfaceId)
    if (selectedApiId.value === `api-${interfaceId}`) {
      closeTab(`api-${interfaceId}`)
    }
    await refreshApiData()
  }

  function reorderApiTreeLocal(folderId: number, orderedApiIds: number[]) {
    const folderNode = findNodeById(apiTree.value, `folder-${folderId}`)
    if (!folderNode?.children) return

    const subFolders = folderNode.children.filter((item) => item.type === 'folder')
    const apiMap = new Map(
      folderNode.children.filter((item) => item.type === 'api').map((item) => [item.id, item]),
    )
    const apis = orderedApiIds
      .map((id) => apiMap.get(`api-${id}`))
      .filter((item): item is ApiTreeNode => !!item)

    if (apis.length !== orderedApiIds.length) return

    folderNode.children = [...subFolders, ...apis]
    apiTree.value = [...apiTree.value]

    const orderMap = new Map(orderedApiIds.map((id, index) => [id, index]))
    interfaces.value = [...interfaces.value].sort((left, right) => {
      if (left.folder_id !== folderId && right.folder_id !== folderId) return 0
      if (left.folder_id !== folderId) return 1
      if (right.folder_id !== folderId) return -1
      return (orderMap.get(left.id) ?? 0) - (orderMap.get(right.id) ?? 0)
    })
  }

  async function submitReorderInterfaces(folderId: number, orderedApiIds: number[]) {
    if (!activeWorkspaceId.value || !activeProjectId.value) return

    reorderApiTreeLocal(folderId, orderedApiIds)
    try {
      await reorderInterfaces(
        activeWorkspaceId.value,
        activeProjectId.value,
        folderId,
        orderedApiIds,
      )
    } catch (error) {
      await refreshApiData()
      throw error
    }
  }

  watch(activeWorkspaceId, () => {
    loadProjects()
  })

  watch(
    [activeModule, activeTabId, workspaceTabs, selectedFolderId, selectedApiId],
    () => {
      persistLayoutState()
    },
    { deep: true },
  )

  watch(activeProjectId, (projectId) => {
    if (!bootstrapComplete || !projectId) return
    refreshApiData()
    loadEnvironments()
  })

  watch(apiTree, () => {
    syncTabLabels()
    if (!navigationHydrated) return

    workspaceTabs.value = workspaceTabs.value.filter((tab) => {
      if (tab.kind === 'create-api') return true
      if (getTabModule(tab) !== 'api') return true
      return !!findNodeById(apiTree.value, tab.id)
    })

    if (activeTabId.value && !workspaceTabs.value.some((tab) => tab.id === activeTabId.value)) {
      const apiTabs = workspaceTabs.value.filter((tab) => getTabModule(tab) === 'api')
      if (apiTabs.length > 0) {
        activateTab(apiTabs[apiTabs.length - 1].id)
      } else {
        clearTabSelection()
        persistLayoutState()
      }
    }
  })

  watch(environments, () => {
    syncTabLabels()

    const hasVariableTabs = workspaceTabs.value.some(
      (tab) => tab.kind === 'env-variables' && tab.id !== 'env-variables-empty',
    )
    if (hasVariableTabs || activeTab.value?.kind === 'env-variables') {
      for (const env of environments.value) {
        upsertTab({
          id: `env-var-${env.id}`,
          kind: 'env-variables',
          module: 'environment',
          label: env.name,
          environmentId: env.id,
          closable: false,
        })
      }
    }

    const envIds = new Set(environments.value.map((item) => item.id))
    workspaceTabs.value = workspaceTabs.value.filter((tab) => {
      if (tab.kind !== 'env-variables' || !tab.environmentId) return true
      return envIds.has(tab.environmentId)
    })

    if (
      activeTab.value?.kind === 'env-variables' &&
      activeTab.value.environmentId &&
      !envIds.has(activeTab.value.environmentId)
    ) {
      if (environments.value.length) {
        openEnvVariableTab(environments.value[0].id, environments.value[0].name)
      } else {
        openEnvSectionTab('variables')
      }
    }
  })

  return {
    activeModule,
    activeWorkspaceId,
    activeProjectId,
    selectedFolderId,
    selectedApiId,
    workspaces,
    projects: workspaceProjects,
    allProjects: projects,
    activeWorkspace,
    activeProject,
    selectedFolder,
    contextMode,
    apiTree,
    interfaces,
    folderInterfaces,
    environments,
    environmentVariables,
    activeEnvironmentId,
    activeEnvironment,
    loadingWorkspaces,
    loadingProjects,
    loadingApiTree,
    loadingInterfaces,
    loadingEnvironments,
    loadingEnvironmentVariables,
    createWorkspaceVisible,
    createProjectVisible,
    creatingApi,
    workspaceTabs,
    moduleTabs,
    activeTabId,
    activeTab,
    activeModuleTab,
    bootstrap,
    refreshWorkspaces,
    refreshProjects,
    refreshApiData,
    refreshEnvironments,
    loadEnvironmentVariables,
    setModule,
    selectWorkspace,
    selectProject,
    selectEnvironment,
    selectFolder,
    selectApi,
    openCreateApi,
    closeCreateApi,
    openProjectSectionTab,
    openSpaceSectionTab,
    openEnvSectionTab,
    openEnvVariableTab,
    openQuickTestTab,
    activateTab,
    closeTab,
    closeTabsLeft,
    closeTabsRight,
    closeOtherTabs,
    closeAllTabs,
    markActiveTabDirty,
    toggleNode,
    isNodeExpanded,
    openCreateWorkspace,
    openCreateProject,
    submitCreateWorkspace,
    submitCreateProject,
    submitCreateFolder,
    submitUpdateFolder,
    submitDeleteFolder,
    submitCreateInterface,
    submitUpdateInterface,
    submitDeleteInterface,
    submitReorderInterfaces,
    parseFolderId,
    parseApiId,
  }
}
