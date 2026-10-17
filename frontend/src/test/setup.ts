import { cleanup } from '@testing-library/react'
import { afterEach, vi } from 'vitest'
import { resetCatalogSession } from '../api/catalogSession.ts'
import '@testing-library/jest-dom/vitest'

afterEach(() => {
  cleanup()
  resetCatalogSession()
  vi.unstubAllEnvs()
  vi.unstubAllGlobals()
})
