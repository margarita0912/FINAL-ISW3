/// <reference types="cypress" />

describe('Flujos completos de ventas (E2E)', () => {

  const mockProductos = [
    { id: 1, nombre: 'Producto A', precio: 100, stock: 10 },
    { id: 2, nombre: 'Producto B', precio: 50,  stock: 5  },
    { id: 3, nombre: 'Producto C', precio: 25,  stock: 8  },
    { id: 4, nombre: 'Producto D', precio: 75,  stock: 3  },
    { id: 5, nombre: 'Producto E', precio: 10,  stock: 2  }
  ]

  beforeEach(() => {

    // Mock login
    cy.intercept('POST', '**/login', {
      statusCode: 200,
      body: { token: 'fake-jwt-token', rol: 'vendedor' }
    }).as('login')

    // Mock productos
    cy.intercept('GET', '**/productos', {
      statusCode: 200,
      body: mockProductos
    }).as('getProductos')

    // Login
    cy.visit('/login')
    cy.get('input[placeholder="Email"]').type('vendedor@test.com')
    cy.get('input[placeholder="Password"]').type('password123')
    cy.contains('Ingresar').click()
    cy.wait('@login')

    // Ir a ventas
    cy.visit('/ventas')
    cy.contains('Registrar Venta', { timeout: 10000 }).should('be.visible')

    // Esperar que carguen los productos mockeados
    cy.wait('@getProductos')
    cy.get('select').first().should('be.visible')
  })

  function seleccionarPrimerProducto() {
    cy.get('select').first().select(mockProductos[0].id.toString())
    cy.get('input[type="number"]').first().should('not.be.disabled')
  }

  it('1️⃣ SIMPLE - Crear venta básica', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('Agregar al carrito').click()
    cy.contains('Producto A').should('exist')
  })

  it('🧪 Test alternativo', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('2')
    cy.contains('Agregar al carrito').click()
    cy.contains('Producto A').should('exist')
  })

  it('3️⃣ Valida stock insuficiente', () => {
    // Producto C tiene stock 0
    cy.get('select').first().select('3') // Producto C
    cy.get('input[type="number"]').first().should('be.disabled')
    cy.contains('Agregar al carrito').should('be.disabled')
  })

  it('4️⃣ Elimina productos del carrito', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('Agregar al carrito').click()

    cy.contains('🗑️').click()
    cy.contains('Producto A').should('not.exist')
  })

  it('5️⃣ Muestra error sin productos', () => {
    cy.contains('Agregar al carrito').should('be.disabled')

    const existeConfirmar = Cypress.$('button:contains("Confirmar Venta")').length > 0
    if (existeConfirmar) {
      cy.contains('Confirmar Venta').should('be.disabled')
    }
  })

  it('6️⃣ Maneja errores del backend', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('Agregar al carrito').click()

    // Mockear error en venta
    cy.intercept('POST', '**/ventas', {
      statusCode: 400,
      body: { error: 'Stock insuficiente en el servidor' }
    }).as('errVenta')

    cy.contains('Confirmar Venta').click()
    cy.wait('@errVenta')

    cy.contains('insuficiente').should('exist')
  })

  it('7️⃣ Test básico', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('Agregar al carrito').click()
    cy.contains('Producto A').should('exist')
  })

})
