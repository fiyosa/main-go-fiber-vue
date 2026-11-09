import { axiosLib } from '../../lib'

export interface LoginRequest {
  email: string
  password: string
}

// LoginResponse mirrors response.LoginResponse from the Go API. The access and
// refresh tokens are delivered as HTTP-only cookies, never in this body.
export interface LoginResponse {
  success: boolean
  message: string
}

/// POST {API_URL}/api/v1/auth/login. The browser stores the issued cookies
/// automatically; no token handling is needed by the caller.
export const login = async (payload: LoginRequest): Promise<LoginResponse> =>
  await axiosLib.instance
    .post<LoginResponse>('/api/v1/auth/login', payload)
    .then((res) => res.data)
    .catch((err) => axiosLib.throwAxios(err))
