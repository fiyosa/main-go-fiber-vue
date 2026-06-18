import { type MaybeRefOrGetter } from 'vue'
import {
  useQuery as vueUseQuery,
  useMutation as vueUseMutation,
  QueryClient,
  type UseQueryOptions,
  type UseQueryReturnType,
  type UseMutationOptions,
  type UseMutationReturnType,
} from '@tanstack/vue-query'

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      retry: false,
      staleTime: 1000 * 60 * 1,
    },
  },
})

export type QueryOptions<TProps = unknown, TResponse = unknown> = Partial<UseQueryOptions<TResponse>> & {
  props?: MaybeRefOrGetter<TProps>
}

export type MutationOptions<TPayload, IResponse> = UseMutationOptions<IResponse, Error, TPayload>

export function useQuery<TProps = unknown, TResponse = unknown>(
  options: QueryOptions<TProps, TResponse>
): UseQueryReturnType<TResponse, Error> {
  return (vueUseQuery as any)(options, queryClient) as UseQueryReturnType<TResponse, Error>
}

export function useMutation<TPayload, IResponse, TOnMutateResult = unknown>(
  options: UseMutationOptions<IResponse, Error, TPayload, TOnMutateResult>
): UseMutationReturnType<IResponse, Error, TPayload, TOnMutateResult> {
  return (vueUseMutation as any)(options, queryClient) as UseMutationReturnType<
    IResponse,
    Error,
    TPayload,
    TOnMutateResult
  >
}
