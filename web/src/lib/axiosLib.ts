import axios, { type AxiosResponse } from 'axios'
import secret from '../config/secret'

export type TypeResponse<T> = AxiosResponse<T>

export const throwAxios = (res: AxiosResponse<any> | any) => {
  if (res?.response) return res.response
  if (res?.request) {
    return {
      status: 400,
      data: {
        message: 'Unable to connect to the server. Check your internet connection.',
      },
    }
  }
  return {
    status: 500,
    data: { message: 'An error occurred that could not be identified.' },
  }
}

export const createQueryStr = (props?: { query?: Record<string, string | number | boolean | undefined> }): string => {
  if (!props?.query) return ''
  const params = new URLSearchParams()
  for (const [key, value] of Object.entries(props.query)) {
    if (value != null && value !== '') params.append(key, String(value))
  }
  const str = params.toString()
  return str ? `?${str}` : ''
}

export const instance = axios.create({
  baseURL: secret.API_URL + '/api',
  timeout: 1000 * 10,
  withCredentials: true,
  headers: {
    Accept: 'application/json',
    'Content-Type': 'application/json',
  },
})
