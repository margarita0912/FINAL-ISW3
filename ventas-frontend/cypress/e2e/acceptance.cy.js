describe('Acceptance - flujo básico de integración', () => {
  it('verifica que la app carga y permite login', () => {
    // Visitar página de login
    cy.visit('/login');
    cy.get('body').should('not.be.empty');
    
    // Hacer login real
    cy.get('input[placeholder="Email"]').type('julio');
    cy.get('input[placeholder="Password"]').type('julio123');
    cy.contains('Ingresar').click();
    
    // Verificar que redirige (no queda en /login)
    cy.url({ timeout: 10000 }).should('not.include', '/login');
    
    // Verificar que el header/nav tiene contenido
    cy.get('body').should('contain', 'Ventas App');
    cy.get('nav').should('exist');
  });
});
