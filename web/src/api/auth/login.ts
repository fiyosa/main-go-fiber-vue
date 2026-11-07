import { axiosLib } from '../../lib'

export interface LoginRequest {
  email: string
  password: string
}

// LoginResponse mirrors response.LoginResponse from the Go API.
export interface LoginResponse {
  success: boolean
  message: string
  token: string
  token_type: string
  expires_in: number
  refresh_token: string
}

/// POST /api/v1/auth/login. The caller is responsible for persisting the token
/// returned by a successful call.
export const login = async (payload: LoginRequest): Promise<LoginResponse> =>
  await axiosLib.instance
    .post<LoginResponse>('/auth/login', payload)
    .then((res) => res.data)
    .catch((err) => axiosLib.throwAxios(err))
