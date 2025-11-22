/// <reference types="cypress" />

export {}

declare global {
  namespace Cypress {
    interface Chainable {
      login(email: string, password: string): Chainable<any>
      addProductToCart(productId: number, quantity: number, discount?: number): Chainable<any>
    }
  }
}

// Comando personalizado para hacer login
Cypress.Commands.add('login', (email: string, password: string) => {
  // Interceptar la llamada de login
  cy.intercept('POST', '**/login', {
    statusCode: 200,
    body: {
      token: 'fake-jwt-token',
      rol: 'vendedor'
    }
  }).as('loginRequest')

  cy.visit('/login')
  cy.get('input[placeholder="Email"]').type(email)
  cy.get('input[placeholder="Password"]').type(password)
  cy.get('button').contains('Ingresar').click()
  
  cy.wait('@loginRequest')
  
  // Verificar que se guardó el token
  cy.window().then((win) => {
    expect(win.localStorage.getItem('token')).to.equal('fake-jwt-token')
    expect(win.localStorage.getItem('rol')).to.equal('vendedor')
  })
})

// Comando para agregar producto al carrito
Cypress.Commands.add('addProductToCart', (productId: number, quantity: number, discount = 0) => {
  // Esperar a que el select esté disponible y tenga opciones
  cy.get('select').eq(0).should('be.visible')
  cy.get('select').eq(0).find('option').should('have.length.greaterThan', 1)
  
  // Intentar seleccionar por valor, si no existe, usar el primer producto disponible
  cy.get('select').eq(0).find(`option[value="${productId}"]`).then($option => {
    if ($option.length > 0) {
      cy.get('select').eq(0).select(productId.toString())
    } else {
      cy.log(`Producto ${productId} no encontrado, seleccionando el primer producto disponible`)
      cy.get('select').eq(0).find('option').eq(1).then($firstOption => {
        cy.get('select').eq(0).select($firstOption.val() as string)
      })
    }
  })
  
  // Esperar a que el botón se habilite después de seleccionar el producto
  cy.get('button').contains('Agregar al carrito').should('not.be.disabled')
  
  // Ingresar cantidad
  cy.get('input[type="number"]').first().clear().type(quantity.toString())
  
  // Aplicar descuento si es necesario
  if (discount > 0) {
    cy.get('input[type="number"]').last().clear().type(discount.toString())
  }
  
  // Hacer click en agregar al carrito
  cy.get('button').contains('Agregar al carrito').click()
})