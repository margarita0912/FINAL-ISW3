import React from 'react'
import { render, screen, waitFor } from '@testing-library/react'
import '@testing-library/jest-dom'
import Productos from '../pages/Productos'

// Mock correcto para axios()
jest.mock('../api/axios', () => {
  return () => Promise.resolve({
    get: jest.fn(() =>
      Promise.resolve({
        data: [
          { id: 1, nombre: 'ProdA', precio: 12.5, stock: 3 }
        ]
      })
    )
  })
})

test('Productos muestra lista obtenida desde API', async () => {
  render(<Productos />)

  await waitFor(() => expect(screen.getByText(/ProdA/)).toBeInTheDocument())

  expect(screen.getByText(/Stock: 3/)).toBeInTheDocument()
})
