import React from 'react'
import { render, screen } from '@testing-library/react'

// Intentamos inicializar MSW; si no es posible, mockeamos axios ANTES de importar el componente
try {
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  require('../mocks/server')
} catch (e) {
  // eslint-disable-next-line global-require
  jest.mock('../api/axios', () => ({
    get: jest.fn(() => Promise.resolve({ data: [{ id: 1, nombre: 'MSW Product A', precio: 12.5, stock: 3 }] })),
  }))
}

import Productos from '../pages/Productos'

// Este test intenta usar MSW si está disponible en el entorno de test.
// En algunos entornos (Jest + transform issues) MSW no puede inicializarse,
// por lo que caemos a un mock de axios para mantener la verificación funcional.
try {
  // Intentamos requerir el server de MSW; si falla, atrapamos y mockeamos axios
  // eslint-disable-next-line @typescript-eslint/no-var-requires
  require('../mocks/server')
} catch (e) {
  // eslint-disable-next-line global-require
  jest.mock('../api/axios', () => ({
    get: jest.fn(() => Promise.resolve({ data: [{ id: 1, nombre: 'MSW Product A', precio: 12.5, stock: 3 }] })),
  }))
}

test('Productos obtiene y muestra productos usando MSW (o fallback)', async () => {
  render(<Productos />)

  // MSW responderá con el producto definido en handlers (si está activo)
  // Aumentamos timeout por si la resolución en entorno CI/Node tarda un poco.
  const item = await screen.findByText(/MSW Product A/, {}, { timeout: 3000 })
  expect(item).toBeInTheDocument()
  expect(screen.getByText(/Stock: 3/)).toBeInTheDocument()
})
