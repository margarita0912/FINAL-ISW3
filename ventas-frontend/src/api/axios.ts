// src/api/axios.ts
import axios from 'axios'

const api = axios.create({
    baseURL: 'http://localhost:8080',
    // withCredentials: true, // dejalo comentado salvo que también uses cookies
})

api.interceptors.request.use((config) => {
    const entorno = localStorage.getItem('entorno') || 'qa'
    const token = localStorage.getItem('token')

    // aseguremos headers definidos (TS-safe)
    config.headers = config.headers || {}

    // tu header de entorno
    ;(config.headers as any)['X-Env'] = entorno


    if (token) {
        ;(config.headers as any)['Authorization'] = `Bearer ${token}`
    }

    return config
})

export default api