import { toValue } from 'vue'
import { tanstackUtil, axiosLib } from '../../lib'

export type LogFileItem = {
  name: string
  size: number
}

interface IProps {
  query?: {
    page?: number
    limit?: number
  }
}

export const getLogs = (options?: tanstackUtil.QueryOptions<IProps, LogFileItem[]>) =>
  tanstackUtil.useQuery<IProps, LogFileItem[]>({
    ...options,
    queryKey: ['logs'],
    queryFn: async () => {
      const p = options?.props ? toValue(options.props) : undefined
      return await axiosLib.instance
        .get('/log' + axiosLib.createQueryStr(p))
        .then((res: axiosLib.TypeResponse<{ data: LogFileItem[] }>) => res.data.data)
        .catch((err) => axiosLib.throwAxios(err))
    },
  })
