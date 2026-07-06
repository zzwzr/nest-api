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
  projectId: number
  name: string
  type: 'folder' | 'api'
  method?: HttpMethod
  children?: ApiTreeNode[]
}
