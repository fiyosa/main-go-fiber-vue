import axios from 'axios'
import type { AxiosInstance, AxiosResponse } from 'axios'

import secret from '../config/secret'

// TypeResponse is the axios response wrapper re-exported so callers do not need
// to import axios directly.
export type TypeResponse<T> = AxiosResponse<T>

export type QueryParams = Record<string, string | number | boolean | null | undefined>

export interface QueryHolder {
  query?: QueryParams
}

// Builds `?a=1&b=2` from an object that may carry a `query` field. Empty
/// values are skipped so optional filters never produce dangling params.
export const createQueryStr = (props?: QueryHolder): string => {
  const params = props?.query
  if (!params) {
    return ''
  }

  const search = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value === undefined || value === null || value === '') {
      continue
    }
    search.set(key, String(value))
  }

  const query = search.toString()
  return query ? `?${query}` : ''
}

/// Normalizes an axios failure into a plain Error with a readable message.
export const throwAxios = (error: unknown): never => {
  if (axios.isAxiosError(error)) {
    const body = error.response?.data as { message?: string } | undefined
    throw new Error(body?.message || error.message)
  }

  throw error instanceof Error ? error : new Error(String(error))
}

export const instance: AxiosInstance = axios.create({
  baseURL: secret.API_URL,
  timeout: 30_000,
  headers: {
    Accept: 'application/json',
  },
})
