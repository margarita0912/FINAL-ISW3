import { useEffect, useState } from 'react'
import api from '../api/axios'

export default function Productos() {
    const [productos, setProductos] = useState([])

    useEffect(() => {
        api.get('/productos')
            .then(res => setProductos(res.data))
            .catch(err => console.error(err))
    }, [])

    return (
        <div>
            <h2>📦 Productos</h2>
            <ul>
                {productos.map((p: any) => (
                    <li key={p.id}>
                        {p.nombre} — ${p.precio} — Stock: {p.stock}
                    </li>
                ))}
            </ul>
        </div>
    )
}