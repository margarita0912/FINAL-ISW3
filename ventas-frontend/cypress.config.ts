import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: process.env.CYPRESS_BASE_URL || 'http://localhost:5174',
    viewportWidth: 1280,
    viewportHeight: 720,
    defaultCommandTimeout: 15000,
    requestTimeout: 15000,
    responseTimeout: 15000,
    pageLoadTimeout: 30000,
    video: false,
    screenshotOnRunFailure: true,

    // No verificar que el servidor esté corriendo antes de los tests
    // Útil cuando apuntamos a servidores remotos
    experimentalStudio: false,

    specPattern: [
      "cypress/e2e/ventas_flow.cy.js",
      "cypress/e2e/acceptance.cy.js"
    ],

    setupNodeEvents(on, config) {
      // Sobrescribir baseUrl desde variable de entorno si existe
      if (process.env.CYPRESS_BASE_URL) {
        config.baseUrl = process.env.CYPRESS_BASE_URL
      }
      return config
    },

    env: {
      apiUrl: process.env.CYPRESS_API_URL || 'http://localhost:8080'
    }
  },
});