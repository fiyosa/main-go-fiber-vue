/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_SHOW_SERVER?: string
  readonly [key: string]: string | undefined
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
