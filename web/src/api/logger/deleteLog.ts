import { tanstackUtil, axiosLib } from '../../lib'

export const deleteLog = (options?: tanstackUtil.MutationOptions<string, any>) =>
  tanstackUtil.useMutation<string, any>({
    mutationFn: async (fileName: string) => await axiosLib.instance.delete(`/log/${encodeURIComponent(fileName)}`),
    ...options,
  })
