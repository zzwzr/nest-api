import http, { getData } from '@/utils/request'
import type { InterfaceItem } from '@/types/workspace'

export function fetchInterfaces(workspaceId: number, projectId: number) {
  return getData<InterfaceItem[]>(
    http.get('/v1/interfaces', { params: { workspace_id: workspaceId, project_id: projectId } }),
  )
}

export function createInterface(
  workspaceId: number,
  projectId: number,
  payload: {
    folder_id: number
    name: string
    method: string
    url?: string
    status?: number
  },
) {
  return getData<{ id: number }>(
    http.post('/v1/interfaces', {
      workspace_id: workspaceId,
      project_id: projectId,
      ...payload,
    }),
  )
}

export function updateInterface(
  workspaceId: number,
  projectId: number,
  interfaceId: number,
  payload: {
    name: string
    method: string
    url?: string
    status?: number
  },
) {
  return getData<null>(
    http.put('/v1/interfaces', {
      workspace_id: workspaceId,
      project_id: projectId,
      interface_id: interfaceId,
      ...payload,
    }),
  )
}

export function deleteInterface(workspaceId: number, projectId: number, interfaceId: number) {
  return getData<null>(
    http.delete('/v1/interfaces', {
      params: {
        workspace_id: workspaceId,
        project_id: projectId,
        interface_id: interfaceId,
      },
    }),
  )
}

export function reorderInterfaces(
  workspaceId: number,
  projectId: number,
  folderId: number,
  interfaceIds: number[],
) {
  return getData<null>(
    http.put('/v1/interfaces/reorder', {
      workspace_id: workspaceId,
      project_id: projectId,
      folder_id: folderId,
      interface_ids: interfaceIds,
    }),
  )
}
