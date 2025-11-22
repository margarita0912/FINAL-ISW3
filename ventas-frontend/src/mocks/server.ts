import { setupServer } from 'msw/node'
import { handlers } from './handlers'

// Setup del server para entorno de tests (Node)
export const server = setupServer(...handlers)
