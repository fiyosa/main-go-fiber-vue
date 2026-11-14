import { axiosLib } from '../../lib'

export interface ChangePasswordRequest {
  'new-password': string
  'confirm-password': string
}

export interface ChangePasswordResponse {
  success: boolean
  message: string
}

/// POST /api/v1/auth/change-password. Success revokes every session of the
/// caller and clears the auth cookies, so the client must sign in again.
/// Throws the raw error so field validation errors (422) can be captured.
export const changePassword = async (payload: ChangePasswordRequest): Promise<ChangePasswordResponse> =>
  await axiosLib.instance
    .post<ChangePasswordResponse>('/api/v1/auth/change-password', payload)
    .then((res) => res.data)
    .catch((err) => {
      throw err
    })
