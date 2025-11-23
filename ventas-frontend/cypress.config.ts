import { defineConfig } from "cypress";

export default defineConfig({
  e2e: {
    baseUrl: 'http://localhost:5174', // QA frontend (Vite runs on 5174)
    viewportWidth: 1280,
    viewportHeight: 720,
    defaultCommandTimeout: 10000,
    requestTimeout: 10000,
    responseTimeout: 10000,
    setupNodeEvents(on, config) {
      // implement node event listeners here
    },
    env: {
      apiUrl: 'http://localhost:8080' // Backend (app) runs on 8080
    }
  },
});
