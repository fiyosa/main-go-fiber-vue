import { axiosLib } from '../../lib'

export interface RegisterRequest {
  email: string
  password: string
}

export interface RegisterResponse {
  success: boolean
  message: string
  data: {
    id: number
    email: string
  }
}

/// POST {API_URL}/api/v1/auth/register. Throws raw error so field validation errors (422) can be captured.
export const register = async (payload: RegisterRequest): Promise<RegisterResponse> =>
  await axiosLib.instance
    .post<RegisterResponse>('/api/v1/auth/register', payload)
    .then((res) => res.data)
    .catch((err) => {
      throw err
    })
