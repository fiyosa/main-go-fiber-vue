import { axiosLib } from '../../lib'

export interface ReadyData {
  status: string
}

interface ReadyResponse {
  success: boolean
  message: string
  data: ReadyData
}

export const getReady = async (): Promise<ReadyData> =>
  await axiosLib.instance
    .get<ReadyResponse>('/api/v1/readyz')
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))
