import { axiosLib } from '../../lib'

export interface LogoutResponse {
  success: boolean
  message: string
}

/// POST /api/v1/auth/logout. Clears auth cookies on the server and revokes session.
export const logout = async (): Promise<LogoutResponse> =>
  await axiosLib.instance
    .post<LogoutResponse>('/api/v1/auth/logout')
    .then((res) => res.data)
    .catch((err) => axiosLib.throwAxios(err))
