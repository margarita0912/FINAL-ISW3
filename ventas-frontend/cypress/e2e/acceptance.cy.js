describe('Acceptance - flujo básico', () => {
  it('verifica que la app carga y el selector de entorno está presente', () => {
    // Usa baseUrl configurado en CI (QA_URL) o local
    cy.visit('/');
    // El selector de entorno muestra un label con el emoji y el texto
    cy.contains('🌐 Entorno actual').should('exist');
    // Aseguramos que el select exista y forzamos entorno QA en localStorage
    cy.window().then((win) => {
      win.localStorage.setItem('entorno', 'qa');
    });
    // Ir a la pantalla de ventas
    cy.visit('/ventas');
    // Comprobamos que la ruta se cargó (no dependemos de contenido específico)
    cy.url().should('include', '/ventas');
    // Comprueba que el body tiene algo de contenido
    cy.get('body').should('not.be.empty');
  });
});
