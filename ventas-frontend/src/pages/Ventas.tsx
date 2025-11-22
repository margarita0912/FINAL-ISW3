import { useEffect, useState } from 'react'
import api from '../api/axios'

type Producto = {
    id: number
    nombre: string
    precio: number
    stock: number
}

type ItemVenta = {
    productoId: number
    nombre: string
    cantidad: number
    precioUnitario: number
    descuento: number
    precioFinal: number
    subtotal: number
}

export default function Ventas() {
    const [productos, setProductos] = useState<Producto[]>([])
    const [productoId, setProductoId] = useState<string>('') // guardamos el ID seleccionado como string
    const [cantidad, setCantidad] = useState(1)
    const [descuento, setDescuento] = useState(0)
    const [items, setItems] = useState<ItemVenta[]>([])
    const [error, setError] = useState('')
    const [filtrarStock, setFiltrarStock] = useState(true)
    const [loading, setLoading] = useState(false)

    useEffect(() => {
        const fetchProductos = async () => {
            try {
                setLoading(true)
                const res = await api.get('/productos')
                // Normalizamos por si el backend usa otras claves o tipos
                const normalizados: Producto[] = (res.data || []).map((raw: any) => ({
                    id: Number(
                        raw.id ??
                        raw.ID ??
                        raw.producto_id ??
                        raw.productoId
                    ),
                    nombre: String(raw.nombre ?? raw.name ?? raw.descripcion ?? raw.descripcion_corta ?? 'Sin nombre'),
                    precio: Number(raw.precio ?? raw.price ?? raw.PRECIO ?? 0),
                    stock: Number(raw.stock ?? raw.STOCK ?? 0),
                }))
                console.log('Productos normalizados:', normalizados)
                setProductos(normalizados)
                setError('')
            } catch (err: any) {
                console.error('Error al cargar productos:', err)
                setError(`Error al cargar productos: ${err.response?.data?.error || err.message}`)
            } finally {
                setLoading(false)
            }
        }
        fetchProductos()
    }, [])

    const agregarItem = () => {
        if (!productoId) {
            setError('Seleccioná un producto válido')
            return
        }

        // Buscamos por igualdad de string para evitar problemas de tipo
        const producto = productos.find(p => String(p.id) === productoId)
        if (!producto) {
            console.warn('IDs disponibles:', productos.map(p => p.id))
            console.warn('productoId seleccionado:', productoId)
            setError('Producto no encontrado')
            return
        }

        if (filtrarStock && producto.stock < cantidad) {
            setError(`Stock insuficiente. Disponible: ${producto.stock}`)
            return
        }

        if (cantidad <= 0) {
            setError('La cantidad debe ser mayor a 0')
            return
        }

        const precioFinal = producto.precio - (producto.precio * descuento / 100)
        const nuevoItem: ItemVenta = {
            productoId: producto.id,
            nombre: producto.nombre,
            cantidad,
            precioUnitario: producto.precio,
            descuento,
            precioFinal,
            subtotal: precioFinal * cantidad
        }

        setItems(prev => [...prev, nuevoItem])
        setProductoId('') // reseteo
        setCantidad(1)
        setDescuento(0)
        setError('')
    }

    const eliminarItem = (index: number) => {
        setItems(items.filter((_, i) => i !== index))
    }

    const registrarVenta = async () => {
        if (items.length === 0) {
            setError('Agregá al menos un producto')
            return
        }

        try {
            setLoading(true)
            setError('')

            const usuarioId = 1 // TODO: reemplazar por el ID real del usuario (auth)

            // Enviamos UNA request por cada item, en el formato que espera tu back
            const operaciones = items.map((item) => {
                const body = {
                    usuario_id: usuarioId,
                    producto_id: item.productoId,
                    cantidad: item.cantidad,
                    descuento: item.descuento
                }
                console.log('Enviando item:', body)
                return api.post('/ventas', body)
            })

            await Promise.all(operaciones)

            alert('✅ Venta registrada exitosamente')
            setItems([])

            // Recargar productos (stock actualizado)
            const resProductos = await api.get('/productos')
            const normalizados: Producto[] = (resProductos.data || []).map((raw: any) => ({
                id: Number(raw.id ?? raw.ID ?? raw.producto_id ?? raw.productoId),
                nombre: String(raw.nombre ?? raw.name ?? raw.descripcion ?? 'Sin nombre'),
                precio: Number(raw.precio ?? raw.price ?? raw.PRECIO ?? 0),
                stock: Number(raw.stock ?? raw.STOCK ?? 0),
            }))
            setProductos(normalizados)
        } catch (err: any) {
            console.error('Error al registrar venta:', err)
            setError(`Error al registrar venta: ${err.response?.data?.error || err.message}`)
        } finally {
            setLoading(false)
        }
    }

    const totalBruto = items.reduce((sum, item) => sum + item.subtotal, 0)
    const iva = totalBruto * 0.21
    const totalFinal = totalBruto + iva

    const productosFiltrados = filtrarStock
        ? productos.filter(p => p.stock > 0)
        : productos

    return (
        <div className="max-w-4xl mx-auto mt-10 p-6 bg-white rounded shadow-md">
            <h2 className="text-2xl font-bold mb-6 text-blue-700">🧾 Registrar Venta</h2>

            {loading && <p className="text-blue-600 mb-4">⏳ Cargando...</p>}
            {error && <p className="text-red-600 mb-4">❌ {error}</p>}

            <div className="mb-4">
                <label className="block font-medium mb-1">Producto</label>
                <select
                    className="border p-2 w-full rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                    value={productoId}
                    onChange={e => {
                        console.log('Select value:', e.target.value)
                        setProductoId(e.target.value)
                        setError('')
                    }}
                    disabled={loading}
                >
                    <option value="">-- Seleccionar producto --</option>
                    {productosFiltrados.map(p => (
                        <option key={p.id} value={String(p.id)}>
                            {p.nombre} - ${p.precio.toFixed(2)} (Stock: {p.stock})
                        </option>
                    ))}
                </select>
                {productos.length === 0 && !loading && (
                    <p className="text-gray-500 text-sm mt-1">No hay productos disponibles</p>
                )}
            </div>

            <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                    <label className="block font-medium mb-1">Cantidad</label>
                    <input
                        type="number"
                        className="border p-2 w-full rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                        value={cantidad}
                        onChange={e => setCantidad(Math.max(1, parseInt(e.target.value) || 1))}
                        min={1}
                        disabled={loading}
                    />
                </div>
                <div>
                    <label className="block font-medium mb-1">Descuento (%)</label>
                    <input
                        type="number"
                        className="border p-2 w-full rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
                        value={descuento}
                        onChange={e => setDescuento(Math.min(100, Math.max(0, parseFloat(e.target.value) || 0)))}
                        min={0}
                        max={100}
                        step={0.1}
                        disabled={loading}
                    />
                </div>
            </div>

            <button
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 w-full mb-4 disabled:bg-gray-400 disabled:cursor-not-allowed"
                onClick={agregarItem}
                disabled={loading || !productoId}
            >
                ➕ Agregar al carrito
            </button>

            <div className="mb-4">
                <label className="inline-flex items-center cursor-pointer">
                    <input
                        type="checkbox"
                        checked={filtrarStock}
                        onChange={e => setFiltrarStock(e.target.checked)}
                        className="mr-2"
                    />
                    <span className="text-sm">Mostrar solo productos con stock</span>
                </label>
            </div>

            {items.length > 0 && (
                <div className="mt-6">
                    <h3 className="text-xl font-semibold mb-2">🛒 Detalle de la venta</h3>
                    <div className="overflow-x-auto">
                        <table className="w-full border text-sm">
                            <thead className="bg-gray-100">
                            <tr>
                                <th className="border px-2 py-2">Producto</th>
                                <th className="border px-2 py-2">Cant.</th>
                                <th className="border px-2 py-2">Precio</th>
                                <th className="border px-2 py-2">Desc.</th>
                                <th className="border px-2 py-2">Subtotal</th>
                                <th className="border px-2 py-2">Acción</th>
                            </tr>
                            </thead>
                            <tbody>
                            {items.map((item, i) => (
                                <tr key={i} className="hover:bg-gray-50">
                                    <td className="border px-2 py-2">{item.nombre}</td>
                                    <td className="border px-2 py-2 text-center">{item.cantidad}</td>
                                    <td className="border px-2 py-2 text-right">${item.precioUnitario.toFixed(2)}</td>
                                    <td className="border px-2 py-2 text-center">{item.descuento}%</td>
                                    <td className="border px-2 py-2 text-right font-semibold">${item.subtotal.toFixed(2)}</td>
                                    <td className="border px-2 py-2 text-center">
                                        <button
                                            onClick={() => eliminarItem(i)}
                                            className="text-red-600 hover:text-red-800 font-bold"
                                            title="Eliminar"
                                        >
                                            🗑️
                                        </button>
                                    </td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>

                    <div className="bg-gray-100 p-4 rounded mt-4">
                        <div className="text-right space-y-1">
                            <p>Total bruto: <strong>${totalBruto.toFixed(2)}</strong></p>
                            <p>IVA (21%): <strong>${iva.toFixed(2)}</strong></p>
                            <p className="text-blue-700 text-xl font-bold border-t-2 border-blue-700 pt-2">
                                Total final: ${totalFinal.toFixed(2)}
                            </p>
                        </div>
                    </div>

                    <button
                        className="bg-green-600 text-white px-6 py-3 rounded hover:bg-green-700 w-full mt-4 font-bold text-lg disabled:bg-gray-400"
                        onClick={registrarVenta}
                        disabled={loading || items.length === 0}
                    >
                        💰 Confirmar Venta
                    </button>
                </div>
            )}
        </div>
    )
}
