import {
  compactFieldTree,
  compactResponseFieldTree,
  emptyFieldNode,
  responseFieldTreeFromApi,
} from '@/utils/interface-field-tree'
import { normalizeResponseExampleForSave } from '@/utils/response-example-format'
import { compactParamRows, emptyParamRow, type ParamRow } from '@/utils/interface-params'
import type {
  HttpMethod,
  HttpProtocol,
  InterfaceDetail,
  InterfaceRequestBody,
  InterfaceResponseExample,
  InterfaceResponseResult,
  InterfaceStatus,
} from '@/types/workspace'

export interface InterfaceEditorFormState {
  protocol: HttpProtocol
  method: HttpMethod
  url: string
  folderId: number | null
  name: string
  status: InterfaceStatus
  requestHeaders: ParamRow[]
  requestBody: InterfaceRequestBody
  queryParams: ParamRow[]
  responseHeaders: ParamRow[]
  responseResults: InterfaceResponseResult[]
  responseExamples: InterfaceResponseExample[]
}

export function createEmptyInterfaceEditorForm(
  defaults?: Partial<Pick<InterfaceEditorFormState, 'method' | 'status'>>,
): InterfaceEditorFormState {
  return {
    protocol: 'HTTP',
    method: defaults?.method ?? 'POST',
    url: '/',
    folderId: null,
    name: '',
    status: defaults?.status ?? 1,
    requestHeaders: [emptyParamRow()],
    requestBody: {
      format: 'json',
      data_type: 'Object',
      raw: '',
      fields: [emptyFieldNode()],
    },
    queryParams: [emptyParamRow()],
    responseHeaders: [emptyParamRow()],
    responseResults: [],
    responseExamples: [],
  }
}

export function defaultResponseResult(name: string): InterfaceResponseResult {
  return {
    name,
    status_code: 200,
    format: 'JSON',
    data_type: 'Object',
    fields: responseFieldTreeFromApi([]),
  }
}

export function defaultResponseExample(name: string): InterfaceResponseExample {
  return normalizeResponseExampleForSave({
    name,
    status_code: 200,
    content_type: 'application/json',
    format: 'JSON',
    data_type: 'Object',
    raw: '',
  })
}

export function buildInterfaceSavePayload(form: InterfaceEditorFormState) {
  return {
    folder_id: form.folderId ?? undefined,
    name: form.name.trim(),
    method: form.method,
    url: form.url.trim(),
    status: form.status,
    request_headers: compactParamRows(form.requestHeaders),
    request_body: {
      format: form.requestBody.format,
      data_type: form.requestBody.data_type,
      raw: form.requestBody.raw ?? '',
      fields: compactFieldTree(form.requestBody.fields),
    },
    query_params: compactParamRows(form.queryParams),
    response_headers: compactParamRows(form.responseHeaders),
    response_results: form.responseResults.map((result) => ({
      ...result,
      fields: compactResponseFieldTree(result.fields),
    })),
    response_examples: form.responseExamples.map((item) => normalizeResponseExampleForSave(item)),
  } satisfies Partial<InterfaceDetail> & {
    folder_id?: number
    name: string
    method: HttpMethod
    url: string
    status: InterfaceStatus
  }
}

/** Stable snapshot for dirty comparison (ignores trailing empty editor rows via compact helpers). */
export function interfaceEditorSnapshot(form: InterfaceEditorFormState) {
  return JSON.stringify({
    protocol: form.protocol,
    ...buildInterfaceSavePayload(form),
  })
}
