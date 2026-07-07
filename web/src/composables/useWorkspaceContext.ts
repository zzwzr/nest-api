import { computed, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createProject as apiCreateProject, fetchProjects } from '@/api/project'
import { createWorkspace as apiCreateWorkspace, fetchWorkspaces } from '@/api/workspace'
import { readWorkspaceLayout, writeWorkspaceLayout } from '@/utils/workspace-layout-storage'
import type { ApiTreeNode, AppModule, ProjectItem, WorkspaceItem } from '@/types/workspace'

const mockApiTree: ApiTreeNode[] = [
  {
    id: 'api-root-1',
    projectId: 0,
    name: '概览',
    type: 'folder',
    children: [
      {
        id: 'api-root-2',
        projectId: 0,
        name: '全部接口',
        type: 'folder',
        children: [
          {
            id: 'api-folder-1',
            projectId: 0,
            name: '用户模块',
            type: 'folder',
            children: [
              {
                id: 'api-folder-2',
                projectId: 0,
                name: '账号相关',
                type: 'folder',
                children: [
                  { id: 'api-1', projectId: 0, name: '退出登录', type: 'api', method: 'POST' },
                  { id: 'api-2', projectId: 0, name: '获取用户信息', type: 'api', method: 'GET' },
                ],
              },
              { id: 'api-3', projectId: 0, name: '修改密码', type: 'api', method: 'PUT' },
            ],
          },
        ],
      },
    ],
  },
]

const storedLayout = readWorkspaceLayout()

const workspaces = ref<WorkspaceItem[]>([])
const projects = ref<ProjectItem[]>([])
const loadingWorkspaces = ref(false)
const loadingProjects = ref(false)

const activeModule = ref<AppModule>('api')
const activeWorkspaceId = ref<number | null>(storedLayout.workspaceId)
const activeProjectId = ref<number | null>(storedLayout.projectId)
const expandedNodeIds = ref<Set<string>>(new Set(['api-root-1', 'api-root-2']))
const selectedApiId = ref<string | null>(null)

const createWorkspaceVisible = ref(false)
const createProjectVisible = ref(false)

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

  const apiTree = computed(() => {
    if (!activeProjectId.value) return []
    return mockApiTree
  })

  async function loadWorkspaces() {
    loadingWorkspaces.value = true
    try {
      const list = await fetchWorkspaces()
      workspaces.value = list

      if (!list.length) {
        activeWorkspaceId.value = null
        activeProjectId.value = null
        projects.value = []
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
        selectedApiId.value = null
      }
    } finally {
      loadingProjects.value = false
    }
  }

  async function bootstrap() {
    await loadWorkspaces()
    await loadProjects()
  }

  async function refreshWorkspaces() {
    await loadWorkspaces()
    await loadProjects()
  }

  async function refreshProjects() {
    await loadProjects()
  }

  function setModule(module: AppModule) {
    activeModule.value = module
  }

  function selectWorkspace(id: number) {
    activeWorkspaceId.value = id
    activeProjectId.value = null
    selectedApiId.value = null
  }

  function selectProject(id: number | null) {
    activeProjectId.value = id
    selectedApiId.value = null
    if (id) {
      expandedNodeIds.value = new Set(['api-root-1', 'api-root-2'])
    }
  }

  function selectApi(id: string) {
    selectedApiId.value = id
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

  watch(activeWorkspaceId, () => {
    loadProjects()
  })

  watch([activeWorkspaceId, activeProjectId], ([workspaceId, projectId]) => {
    writeWorkspaceLayout({ workspaceId, projectId })
  })

  return {
    activeModule,
    activeWorkspaceId,
    activeProjectId,
    selectedApiId,
    workspaces,
    projects: workspaceProjects,
    allProjects: projects,
    activeWorkspace,
    activeProject,
    contextMode,
    apiTree,
    loadingWorkspaces,
    loadingProjects,
    createWorkspaceVisible,
    createProjectVisible,
    bootstrap,
    refreshWorkspaces,
    refreshProjects,
    setModule,
    selectWorkspace,
    selectProject,
    selectApi,
    toggleNode,
    isNodeExpanded,
    openCreateWorkspace,
    openCreateProject,
    submitCreateWorkspace,
    submitCreateProject,
  }
}
