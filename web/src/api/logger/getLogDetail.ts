import { computed, toValue } from 'vue'
import { tanstackUtil, axiosLib } from '../../lib'
import type { IResPaginate } from '../base'

type LogEntry = {
  level: string
  time: string
  message: string
}

type LogDetailResponse = IResPaginate<LogEntry[]>

interface IProps {
  param: {
    file_name: string
  }
  query?: {
    search?: string
    page?: string
    limit?: string
    levels?: string
  }
}

export const getLogDetail = (options?: tanstackUtil.QueryOptions<IProps, LogDetailResponse>) => {
  const resolved = computed(() => (options?.props ? toValue(options.props) : undefined))
  return tanstackUtil.useQuery<IProps, LogDetailResponse>({
    ...options,
    queryKey: computed(() => [
      'log',
      'detail',
      resolved.value?.param?.file_name,
      resolved.value?.query?.search,
      resolved.value?.query?.page,
      resolved.value?.query?.levels,
    ]),
    queryFn: async () => {
      const p = resolved.value
      if (!p) return { message: '', data: [], meta: { total: 0, page: 1, limit: 50 } }
      return await axiosLib.instance
        .get(`/log/` + encodeURIComponent(p.param.file_name) + axiosLib.createQueryStr(p))
        .then((res: axiosLib.TypeResponse<LogDetailResponse>) => res.data)
        .catch((err) => axiosLib.throwAxios(err))
    },
    enabled: computed(() => !!resolved.value?.param?.file_name),
  })
}
