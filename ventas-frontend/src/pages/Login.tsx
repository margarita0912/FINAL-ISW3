import { useState } from 'react'
import api from '../api/axios'

export default function Login() {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState('')

    const handleLogin = async () => {
        try {
            const res = await api.post('/login', { nombre: email, clave: password })
            localStorage.setItem('token', res.data.token)
            localStorage.setItem('rol', res.data.rol)
            window.location.href = '/'
        } catch (err) {
            setError('Credenciales inválidas')
        }
    }

    return (
        <div>
            <h2>🔐 Login</h2>
            <input placeholder="Email" value={email} onChange={e => setEmail(e.target.value)} />
            <input placeholder="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} />
            <button onClick={handleLogin}>Ingresar</button>
            {error && <p style={{ color: 'red' }}>{error}</p>}
        </div>
    )
}
