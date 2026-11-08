import { axiosLib } from '../../lib'

export interface HealthData {
  status: string
}

interface HealthResponse {
  success: boolean
  message: string
  data: HealthData
}

/// GET {API_URL}/api/v1/healthz. Useful as a smoke test for the axios setup.
export const getHealth = async (): Promise<HealthData> =>
  await axiosLib.instance
    .get<HealthResponse>('/api/v1/healthz')
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))
