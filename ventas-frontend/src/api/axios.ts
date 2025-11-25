import axios from 'axios'

// 1) Obtener la URL del entorno (Render la pone en build)
const envURL = import.meta.env.VITE_API_URL

// 2) Elegimos qué base usar
// Render: usa envURL
// Cypress/local: usa localhost
const baseURL = envURL && envURL.trim() !== ""
  ? envURL
  : "http://localhost:8080"

// 3) Crear instancia axios
const api = axios.create({
  baseURL
})

// 4) Token (normal)
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  config.headers = config.headers || {}
  if (token) config.headers["Authorization"] = `Bearer ${token}`
  return config
})

console.log("API URL usada:", import.meta.env.VITE_API_URL)


export default api
