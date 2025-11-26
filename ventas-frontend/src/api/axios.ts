import axios, { AxiosInstance } from "axios";

let apiInstance: AxiosInstance | null = null;
let loadingPromise: Promise<void> | null = null;

// Cargar config.json una sola vez
async function loadConfigJson() {
  try {
    const res = await fetch("/config.json", { cache: "no-store" });
    const json = await res.json();
    // Permitir string vacío "" para usar base relativa en producción
    if (Object.prototype.hasOwnProperty.call(json, "api_url")) {
      return json.api_url as string;
    }
    return null;
  } catch (err) {
    console.warn("No se pudo cargar config.json, usando fallback.");
    return null;
  }
}

// Inicializar Axios (solo una vez)
async function init() {
  if (apiInstance) return; // ya inicializado

  const runtimeURL = await loadConfigJson();
  let finalURL: string | undefined;
  if (runtimeURL !== null && runtimeURL !== undefined) {
    finalURL = runtimeURL; // puede ser "" (base relativa)
  } else if (typeof import.meta.env.VITE_API_URL !== "undefined") {
    finalURL = import.meta.env.VITE_API_URL as string; // puede ser ""
  } else {
    finalURL = "http://localhost:8080";
  }

  console.log("API URL usada:", finalURL);

  apiInstance = axios.create({
    baseURL: finalURL,
  });

  apiInstance.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    config.headers = config.headers || {};
    if (token) config.headers["Authorization"] = `Bearer ${token}`;
    return config;
  });
}

// API que devuelve SIEMPRE un AxiosInstance válido
async function getAPI(): Promise<AxiosInstance> {
  if (!loadingPromise && !apiInstance) {
    loadingPromise = init();
  }
  if (loadingPromise) {
    await loadingPromise;
    loadingPromise = null;
  }
  return apiInstance!;
}

export default getAPI;
