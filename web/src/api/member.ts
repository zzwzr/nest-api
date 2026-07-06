import http, { getData } from '@/utils/request'
import type { MemberItem, WorkspaceRole } from '@/types/workspace'

export function fetchMembers(workspaceId: number) {
  return getData<MemberItem[]>(
    http.get('/v1/members', { params: { workspace_id: workspaceId } }),
  )
}

export function inviteMember(workspaceId: number, userId: number, role: WorkspaceRole) {
  return getData<null>(
    http.post('/v1/members', { workspace_id: workspaceId, user_id: userId, role }),
  )
}

export function updateMemberRole(workspaceId: number, memberId: number, role: WorkspaceRole) {
  return getData<null>(
    http.put('/v1/members', { workspace_id: workspaceId, member_id: memberId, role }),
  )
}

export function removeMember(workspaceId: number, memberId: number) {
  return getData<null>(
    http.delete('/v1/members', { params: { workspace_id: workspaceId, member_id: memberId } }),
  )
}
