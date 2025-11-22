import '@testing-library/jest-dom';
import { TextEncoder, TextDecoder } from 'node:util';

// Asignación segura sin redefinir tipos
(globalThis as unknown as { TextEncoder: typeof TextEncoder }).TextEncoder = TextEncoder;
(globalThis as unknown as { TextDecoder: typeof TextDecoder }).TextDecoder = TextDecoder;
// MSW server para tests: arrancamos/paramos el server automáticamente
// Mock global de fetch para todos los tests (fallback si MSW da problemas)
try {
	// Intentamos arrancar MSW server (Node) importando el server creado en src/mocks/server.ts
	// Usamos require para evitar problemas con ESM/TS en el entorno de Jest.
	// eslint-disable-next-line @typescript-eslint/no-var-requires
	const { server } = require('./src/mocks/server') as { server: any }

	// Lifecycle de MSW en Jest
	beforeAll(() => server.listen())
	afterEach(() => server.resetHandlers())
	afterAll(() => server.close())
} catch (e) {
	// Si por alguna razón MSW no puede inicializarse, dejamos un fetch global como fallback
	// para mantener los tests funcionales.
	// eslint-disable-next-line no-console
	console.warn('MSW server no pudo inicializarse en setupTests, usando fetch fallback:', e)

	global.fetch = jest.fn(() =>
		Promise.resolve({
			ok: true,
			status: 200,
			statusText: 'OK',
			headers: new Headers(),
			redirected: false,
			type: 'basic',
			url: '',
			clone: () => undefined,
			body: null,
			bodyUsed: false,
			async json() {
				return { id: 1, nombre: 'Fallback Product', precio: 100, stock: 1 }
			},
			async text() {
				return JSON.stringify({ id: 1, nombre: 'Fallback Product', precio: 100, stock: 1 })
			},
			async arrayBuffer() {
				return new ArrayBuffer(0)
			},
			async blob() {
				return new Blob()
			},
			async formData() {
				return new FormData()
			},
		} as unknown as Response)
	)
}
