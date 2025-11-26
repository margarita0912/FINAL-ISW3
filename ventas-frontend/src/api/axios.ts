import axios, { AxiosInstance } from "axios";

let api: AxiosInstance | null = null;

// URL por defecto
let finalURL = import.meta.env.VITE_API_URL || "http://localhost:8080";

// Cargar config.json en runtime
async function loadRuntimeConfig() {
  try {
    const res = await fetch("/config.json", { cache: "no-store" });

    if (res.ok) {
      const json = await res.json();
      if (json.api_url) {
        finalURL = json.api_url;
      }
    } else {
      console.warn("config.json no encontrado, usando fallback.");
    }
  } catch (err) {
    console.warn("No se pudo cargar config.json, usando fallback.");
  }
}

// Inicializar API (se llama una sola vez)
async function initAPI() {
  if (api) return api; // ya inicializado

  await loadRuntimeConfig();

  console.log("API URL usada:", finalURL);

  api = axios.create({
    baseURL: finalURL,
  });

  // Interceptor de token
  api.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) {
      config.headers = config.headers || {};
      config.headers["Authorization"] = `Bearer ${token}`;
    }
    return config;
  });

  return api;
}

// Exportamos la función (no la instancia directa)
export default initAPI;
