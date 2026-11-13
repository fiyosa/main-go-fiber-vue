import { axiosLib } from '../../lib'

export interface UserMe {
  id: number
  email: string
  roles: string[]
}

interface MeResponse {
  success: boolean
  message: string
  data: UserMe
}

/// GET {API_URL}/api/v1/auth/me. Fetches the authenticated user identity via HTTP-only cookie.
export const getMe = async (): Promise<UserMe> =>
  await axiosLib.instance
    .get<MeResponse>('/api/v1/auth/me')
    .then((res) => res.data.data)
    .catch((err) => axiosLib.throwAxios(err))
