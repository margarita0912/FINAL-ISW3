process.env.MSW_DISABLE_SOCKET = 'true'; // ← esto evita BroadcastChannel
require('whatwg-fetch');
const { TextEncoder, TextDecoder } = require('util');
global.TextEncoder = TextEncoder;
global.TextDecoder = TextDecoder;

// Mock de import.meta.env para Vite en Jest
global.import = {
	meta: {
		env: {
			VITE_API_URL: 'http://localhost:8080',
			MODE: 'test'
		}
	}
};

// Polyfill TransformStream en el entorno de Jest si no está definido.
// Preferimos la implementación nativa de Node (`stream/web`) cuando esté disponible.
try {
	if (typeof global.TransformStream === 'undefined') {
		// Node.js (>= 16.7) expone TransformStream en 'stream/web'
		// Usamos require para evitar errores en entornos antiguos.
		// eslint-disable-next-line global-require
		const streamWeb = require('stream/web');
		if (streamWeb && streamWeb.TransformStream) {
			global.TransformStream = streamWeb.TransformStream;
		}
	}
} catch (e) {
	// noop: si no está disponible, MSW necesitará otro polyfill; lo dejamos pasar
}
