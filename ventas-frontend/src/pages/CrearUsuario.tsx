import { useState, useEffect } from 'react'
import api from '../api/axios'

export default function CrearUsuario() {
    const [nombre, setNombre] = useState('')
    const [clave, setClave] = useState('')
    const [rolNuevo, setRolNuevo] = useState('comprador')
    const [error, setError] = useState('')
    const [exito, setExito] = useState('')
    const [permitido, setPermitido] = useState(false)

    useEffect(() => {
        const rol = localStorage.getItem('rol')
        setPermitido(rol === 'precio')
    }, [])

    const crear = async () => {
        setError('')
        setExito('')
        try {
            const token = localStorage.getItem('token')
            await api.post('/usuarios', { nombre, clave, rol: rolNuevo }, {
                headers: { Authorization: `Bearer ${token}` }
            })
            setExito('✅ Usuario creado con éxito')
            setNombre('')
            setClave('')
            setRolNuevo('comprador')
        } catch (err) {
            setError('❌ Error al crear usuario')
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
            <h2 className="text-2xl font-bold mb-4 text-blue-700">👤 Crear Usuario</h2>

            <label className="block mb-2 font-medium">Nombre</label>
            <input
                className="border p-2 w-full mb-4 rounded"
                value={nombre}
                onChange={e => setNombre(e.target.value)}
                placeholder="Ej: lucas"
            />

            <label className="block mb-2 font-medium">Clave</label>
            <input
                className="border p-2 w-full mb-4 rounded"
                type="password"
                value={clave}
                onChange={e => setClave(e.target.value)}
                placeholder="Ej: clave123"
            />

            <label className="block mb-2 font-medium">Rol</label>
            <select
                className="border p-2 w-full mb-4 rounded"
                value={rolNuevo}
                onChange={e => setRolNuevo(e.target.value)}
            >
                <option value="comprador">Comprador</option>
                <option value="vendedor">Vendedor</option>
                <option value="precio">Precio</option>
            </select>

            <button
                className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 w-full"
                onClick={crear}
            >
                Crear Usuario
            </button>

            {exito && <p className="text-green-600 mt-4">{exito}</p>}
            {error && <p className="text-red-600 mt-4">{error}</p>}
        </div>
    )
}
