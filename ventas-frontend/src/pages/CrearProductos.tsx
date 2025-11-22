import { useState, useEffect } from 'react'
import api from '../api/axios'

export default function CrearProductos() {
    const [nombre, setNombre] = useState('')
    const [precio, setPrecio] = useState<number | ''>('')
    const [stock, setStock] = useState<number | ''>('')
    const [error, setError] = useState('')
    const [exito, setExito] = useState('')
    const [permitido, setPermitido] = useState(false)

    useEffect(() => {
        const rol = localStorage.getItem('rol')
        setPermitido(rol === 'vendedor' || rol === 'comprador')
    }, [])

    const crear = async () => {
        setError('')
        setExito('')
        if (!nombre || precio === '' || stock === '') {
            setError('Todos los campos son obligatorios')
            return
        }

        try {
            const token = localStorage.getItem('token')
            await api.post('/productos', { nombre, precio, stock }, {
                headers: { Authorization: `Bearer ${token}` }
            })
            setExito('✅ Producto creado con éxito')
            setNombre('')
            setPrecio('')
            setStock('')
        } catch (err) {
            setError('❌ Error al crear producto')
        }
    }

    if (!permitido) {
        return (
            <div className="text-center mt-10 text-red-600 font-semibold">
                ⛔ No tenés permisos para acceder a esta sección.
            </div>
        )
    }

    return (
        <div className="max-w-md mx-auto mt-10 bg-white p-6 rounded shadow-md border border-gray-200">
            <h2 className="text-2xl font-bold mb-4 text-blue-700">📦 Crear Producto</h2>

            <label className="block mb-2 font-medium">Nombre</label>
            <input
                className="border p-2 w-full mb-4 rounded"
                value={nombre}
                onChange={e => setNombre(e.target.value)}
                placeholder="Ej: Camisa"
            />

            <label className="block mb-2 font-medium">Precio</label>
            <input
                className="border p-2 w-full mb-4 rounded"
                type="number"
                value={precio}
                onChange={e => setPrecio(parseFloat(e.target.value))}
                placeholder="Ej: 2500"
                min={0}
            />

            <label className="block mb-2 font-medium">Stock</label>
            <input
                className="border p-2 w-full mb-4 rounded"
                type="number"
                value={stock}
                onChange={e => setStock(parseInt(e.target.value))}
                placeholder="Ej: 10"
                min={0}
            />

            <button
                className="bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 w-full"
                onClick={crear}
            >
                Crear Producto
            </button>

            {exito && <p className="text-green-600 mt-4">{exito}</p>}
            {error && <p className="text-red-600 mt-4">{error}</p>}
        </div>
    )
}
