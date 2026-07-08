export type AppModule = 'api' | 'quick-test' | 'environment' | 'project'

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS'

export type WorkspaceRole = 1 | 2 | 3 | 4

export interface WorkspaceItem {
  id: number
  name: string
  owner_id: number
  owner_name?: string
  role: WorkspaceRole
  created_at: string
}

export interface ProjectItem {
  id: number
  workspace_id: number
  name: string
  created_by: number
  creator_name?: string
  created_at: string
}

export interface MemberItem {
  id: number
  user_id: number
  name: string
  account: string
  avatar: string
  role: WorkspaceRole
  created_at: string
}

export interface ApiTreeNode {
  id: string
  project_id: number
  name: string
  type: 'folder' | 'api'
  method?: HttpMethod
  children?: ApiTreeNode[]
}

export interface InterfaceItem {
  id: number
  project_id: number
  folder_id: number
  name: string
  method: HttpMethod
  url: string
  status: 1 | 2
  folder_name?: string
  updated_by?: number
  updated_by_name?: string
  created_at: string
  updated_at?: string
}

export type WorkspaceTabKind =
  | 'folder'
  | 'api'
  | 'create-api'
  | 'workspace-list'
  | 'project-list'
  | 'member-list'
  | 'env-list'
  | 'env-variables'
  | 'quick-test'

export interface WorkspaceTab {
  id: string
  kind: WorkspaceTabKind
  module: AppModule
  label: string
  method?: HttpMethod
  folderId?: string
  apiId?: string
  environmentId?: number
  dirty?: boolean
  closable?: boolean
}

export interface EnvironmentVariableItem {
  id: number
  environment_id: number
  key: string
  value: string
  description: string
  created_at: string
  updated_at: string
}

export interface EnvironmentItem {
  id: number
  project_id: number
  name: string
  base_url: string
  is_default: boolean
  created_at: string
}
