export type IRes = {
  message: string
}

export type IError = {
  message: string
  errors: Record<string, object[]> | null
}

export type IResPaginate<TData = unknown> = {
  message: string
  data: TData
  meta: IMeta
}

export type IMeta = {
  total: number
  page: number
  limit: number
}
