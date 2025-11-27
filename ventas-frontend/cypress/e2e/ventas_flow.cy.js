/// <reference types="cypress" />

describe('Flujos completos de ventas (E2E - Integración QA)', () => {

  beforeEach(() => {
    // Login REAL contra QA
    cy.visit('/login')
    cy.get('input[placeholder="Email"]', { timeout: 10000 }).should('be.visible').type('ariel')
    cy.get('input[placeholder="Password"]').type('piqui123')
    cy.contains('Ingresar').click()
    
    // Esperar redirección tras login exitoso (puede ir a /productos o /ventas)
    cy.url({ timeout: 15000 }).should('not.include', '/login')

    // Ir a ventas directamente
    cy.visit('/ventas')
    cy.contains('Registrar Venta', { timeout: 15000 }).should('be.visible')

    // Esperar a que carguen productos reales del select
    cy.get('select').first({ timeout: 10000 }).should('be.visible')
    cy.get('select option', { timeout: 5000 }).should('have.length.gt', 1)
  })

  function seleccionarPrimerProducto() {
    // Seleccionar el primer producto disponible (índice 1, ya que 0 es "Selecciona...")
    cy.get('select').first().find('option').eq(1).then($option => {
      cy.get('select').first().select($option.val())
    })
    cy.get('input[type="number"]').first().should('not.be.disabled')
  }

  it('1️⃣ SIMPLE - Crear venta básica', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('button', 'Agregar al carrito').should('not.be.disabled').click()
    
    // Wait y verificar que aparece en el carrito
    cy.wait(1500)
    cy.get('body').should('contain', 'Detalle de la venta')
    cy.get('table tbody tr').should('have.length', 1)
  })

  it('🧪 Test alternativo - cantidad 2', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('2')
    cy.contains('button', 'Agregar al carrito').should('not.be.disabled').click()
    cy.wait(1500)
    cy.get('body').should('contain', 'Detalle de la venta')
    cy.get('table tbody tr').should('have.length', 1)
  })

  it('3️⃣ Valida stock insuficiente', () => {
    // Seleccionar segundo producto disponible
    cy.get('select').first().find('option').eq(2).then($option => {
      if ($option.length) {
        cy.get('select').first().select($option.val())
        cy.get('input[type="number"]').first().clear().type('9999')
        cy.contains('Agregar al carrito').click()
        cy.contains(/stock insuficiente|insuficiente/i).should('exist')
      } else {
        cy.log('Skip: no hay suficientes productos para este test')
      }
    })
  })

  it('4️⃣ Elimina productos del carrito', () => {
    cy.reload()
    cy.contains('Registrar Venta', { timeout: 10000 }).should('be.visible')

    // Asegurar que el checkbox de filtrar stock esté desmarcado
    cy.get('input[type="checkbox"]').uncheck()

    // Agregar primer producto disponible
    cy.get('select').first().find('option').eq(1).then($option1 => {
      const nombre1 = $option1.text().split(' -')[0] // Extraer solo el nombre sin precio/stock
      cy.get('select').first().select($option1.val())
      cy.get('input[type="number"]').first().clear().type('1')
      cy.contains('button', 'Agregar al carrito').should('not.be.disabled').click()
      cy.contains('Detalle de la venta').should('be.visible')
      cy.get('table tbody').should('contain', nombre1)

      // Agregar segundo producto
      cy.get('select').first().find('option').eq(2).then($option2 => {
        const nombre2 = $option2.text().split(' -')[0]
        cy.get('select').first().select($option2.val())
        cy.get('input[type="number"]').first().clear().type('1')
        cy.contains('button', 'Agregar al carrito').should('not.be.disabled').click()
        cy.get('table tbody').should('contain', nombre2)
        cy.get('table tbody tr').should('have.length', 2)

        // Eliminar el primero (encontrar el botón de eliminar en la primera fila)
        cy.get('table tbody tr').first().find('button[title="Eliminar"]').click()

        // Verificar que ahora solo queda 1 fila
        cy.get('table tbody tr').should('have.length', 1)
        cy.get('table tbody').should('contain', nombre2)
      })
    })
  })

  it('5️⃣ Muestra error sin productos', () => {
    cy.contains('Agregar al carrito').should('be.disabled')

    // Verificar que no existe botón de confirmar o está deshabilitado
    cy.get('body').then($body => {
      if ($body.find('button').filter((_, el) => /confirmar/i.test(el.textContent)).length > 0) {
        cy.contains('button', /confirmar/i).should('be.disabled')
      }
    })
  })

  it('6️⃣ Confirma venta exitosa', () => {
    cy.reload()
    cy.contains('Registrar Venta', { timeout: 10000 }).should('be.visible')

    // Asegurar que el checkbox de filtrar stock esté desmarcado
    cy.get('input[type="checkbox"]').uncheck()

    // Agregar un producto
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('button', 'Agregar al carrito').should('not.be.disabled').click()

    // Esperar que aparezca en el carrito
    cy.wait(1500)
    cy.get('body').should('contain', 'Detalle de la venta')
    cy.get('table tbody tr').should('have.length', 1)

    // Confirmar venta
    cy.contains('button', /confirmar/i).should('be.visible').and('not.be.disabled').click()

    // Esperar respuesta del servidor y verificar que el carrito se limpió
    cy.wait(3000)
    cy.get('body').should('not.contain', 'Detalle de la venta')
  })

  it('7️⃣ Test básico - agregar producto', () => {
    seleccionarPrimerProducto()
    cy.get('input[type="number"]').first().clear().type('1')
    cy.contains('button', 'Agregar al carrito').should('not.be.disabled').click()
    cy.wait(1500)
    cy.get('body').should('contain', 'Detalle de la venta')
    cy.get('table tbody tr').should('have.length', 1)
  })

})
