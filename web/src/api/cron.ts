import { axiosLib } from '../lib'

export interface CronEntity {
  id: string
  name: string
  schedule: string
  timezone: string
  script: string
  env?: Record<string, string>
  overlap_policy: 'skip' | 'allow'
  enabled: boolean
  created_at: string
  updated_at: string
}

export interface CronRunEntity {
  id: string
  cron_id: string
  status: string
  trigger_type: string
  started_at: string
  finished_at?: string
  duration_ms: number
  output?: string
  error?: string
  created_at: string
}

export interface CronCreateRequest {
  name: string
  schedule: string
  timezone: string
  script: string
  env?: Record<string, string>
  overlap_policy: 'skip' | 'allow'
  enabled?: boolean
}

export interface CronUpdateRequest {
  name: string
  schedule: string
  timezone: string
  script: string
  env?: Record<string, string>
  overlap_policy: 'skip' | 'allow'
  enabled?: boolean
}

export interface ApiResponse<T> {
  success: boolean
  message: string
  data: T
}

export const listCrons = async (): Promise<CronEntity[]> =>
  await axiosLib.instance
    .get<ApiResponse<CronEntity[]>>('/api/v1/crons')
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const getCron = async (id: string): Promise<CronEntity> =>
  await axiosLib.instance
    .get<ApiResponse<CronEntity>>(`/api/v1/crons/${id}`)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const createCron = async (payload: CronCreateRequest): Promise<CronEntity> =>
  await axiosLib.instance
    .post<ApiResponse<CronEntity>>('/api/v1/crons', payload)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const updateCron = async (id: string, payload: CronUpdateRequest): Promise<CronEntity> =>
  await axiosLib.instance
    .put<ApiResponse<CronEntity>>(`/api/v1/crons/${id}`, payload)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const deleteCron = async (id: string): Promise<void> =>
  await axiosLib.instance
    .delete<ApiResponse<null>>(`/api/v1/crons/${id}`)
    .then(() => {})
    .catch((err) => axiosLib.throwAxios(err))

export const toggleCron = async (id: string): Promise<CronEntity> =>
  await axiosLib.instance
    .post<ApiResponse<CronEntity>>(`/api/v1/crons/${id}/toggle`)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const triggerCron = async (id: string): Promise<CronRunEntity> =>
  await axiosLib.instance
    .post<ApiResponse<CronRunEntity>>(`/api/v1/crons/${id}/trigger`)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const listCronRuns = async (id: string): Promise<CronRunEntity[]> =>
  await axiosLib.instance
    .get<ApiResponse<CronRunEntity[]>>(`/api/v1/crons/${id}/runs`)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const listAllRuns = async (): Promise<CronRunEntity[]> =>
  await axiosLib.instance
    .get<ApiResponse<CronRunEntity[]>>('/api/v1/crons/runs')
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))

export const getRunDetail = async (runId: string): Promise<CronRunEntity> =>
  await axiosLib.instance
    .get<ApiResponse<CronRunEntity>>(`/api/v1/crons/runs/${runId}`)
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))
