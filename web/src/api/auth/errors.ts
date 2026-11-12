import axios from 'axios'

// AuthError is the UI-friendly shape of a failed auth request. The Go API
// answers 422 with per-field messages, every other status with a single
// translated message; both are folded into this structure.
export interface AuthError {
  message: string
  /** Field name -> translated validation messages. */
  fields: Record<string, string[]>
}

const FALLBACK_MESSAGE = 'Unable to sign in. Please try again.'
const NETWORK_MESSAGE = 'Cannot reach the server. Check your connection and try again.'

interface ErrorBody {
  message?: string
  errors?: Record<string, string[]>
}

/// Converts any thrown value into an AuthError without losing field details.
/// throwAxios is not used here because it collapses the body to one string.
export const toAuthError = (error: unknown): AuthError => {
  if (axios.isAxiosError(error)) {
    const body = error.response?.data as ErrorBody | undefined

    if (!error.response) {
      return { message: NETWORK_MESSAGE, fields: {} }
    }

    return {
      message: body?.message || FALLBACK_MESSAGE,
      fields: body?.errors ?? {},
    }
  }

  if (error instanceof Error && error.message) {
    return { message: error.message, fields: {} }
  }

  return { message: FALLBACK_MESSAGE, fields: {} }
}

/// Reads the first message for a field so templates do not repeat the lookup.
export const firstFieldError = (fields: Record<string, string[]>, name: string): string => fields[name]?.[0] ?? ''
