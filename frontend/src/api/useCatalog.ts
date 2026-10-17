import { useSyncExternalStore } from 'react'
import {
  getCatalogSnapshot,
  subscribeCatalog,
  type CatalogSnapshot,
} from './catalogSession.ts'

export function useCatalog(): CatalogSnapshot {
  return useSyncExternalStore(
    subscribeCatalog,
    getCatalogSnapshot,
    getCatalogSnapshot,
  )
}
