// Usamos require dinámico para evitar errores de tipos/ESM entre distintas versiones
// de `msw` en el entorno de tests. Esto evita errores de compilación relacionados
// con las exportaciones de la librería.
// eslint-disable-next-line @typescript-eslint/no-var-requires
const rest = (require('msw') as any).rest

// Utilidad: productos de ejemplo que pueden compartirse entre handlers
const sampleProducts = [
  { id: 1, nombre: 'MSW Product A', precio: 12.5, stock: 3 },
]

// Exportar handlers. Cubrimos tanto rutas absolutas como relativas para evitar
// que las peticiones no sean interceptadas por diferencias en la baseURL.
export const handlers = [
  // Handler genérico que atrapa cualquier URL que contenga 'productos'
  rest.get(/productos/, (req: any, res: any, ctx: any) => {
    return res(ctx.status(200), ctx.json(sampleProducts))
  }),
  // GET productos (ruta relativa)
  rest.get('/productos', (req: any, res: any, ctx: any) => {
    return res(ctx.status(200), ctx.json(sampleProducts))
  }),

  // GET productos (ruta absoluta típica con axios baseURL)
  rest.get('http://localhost:8080/productos', (req: any, res: any, ctx: any) => {
    return res(ctx.status(200), ctx.json(sampleProducts))
  }),

  // POST login: return 401 para credenciales vacías, 200 con token para otros
  rest.post('/login', async (req: any, res: any, ctx: any) => {
    const body = await req.json()
    if (!body || !body.nombre || !body.clave) {
      return res(ctx.status(401), ctx.json({ error: 'Datos inválidos' }))
    }
    return res(ctx.status(200), ctx.json({ token: 'fake-jwt', rol: 'user' }))
  }),

  rest.post('http://localhost:8080/login', async (req: any, res: any, ctx: any) => {
    const body = await req.json()
    if (!body || !body.nombre || !body.clave) {
      return res(ctx.status(401), ctx.json({ error: 'Datos inválidos' }))
    }
    return res(ctx.status(200), ctx.json({ token: 'fake-jwt', rol: 'user' }))
  }),

  // POST compras/ventas (simulaciones básicas)
  rest.post('/compras', async (req: any, res: any, ctx: any) => {
    const body = await req.json()
    if (!body || !body.producto_id || !body.cantidad || body.cantidad <= 0) {
      return res(ctx.status(400), ctx.json({ error: 'Datos inválidos' }))
    }
    return res(ctx.status(201), ctx.json({ ok: true }))
  }),

  rest.post('/ventas', async (req: any, res: any, ctx: any) => {
    const body = await req.json()
    if (!body || !body.producto_id || !body.cantidad || body.cantidad <= 0) {
      return res(ctx.status(400), ctx.json({ error: 'Datos inválidos' }))
    }
    // calcular precio de ejemplo
    const prod = sampleProducts.find((p) => p.id === body.producto_id) || sampleProducts[0]
    const precio = prod.precio * body.cantidad
    return res(ctx.status(201), ctx.json({ precio_final: precio }))
  }),
]
