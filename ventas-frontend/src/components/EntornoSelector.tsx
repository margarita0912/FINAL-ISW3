import { useState } from 'react'

export default function EntornoSelector() {
    const [entorno, setEntorno] = useState(localStorage.getItem('entorno') || 'qa')

    const cambiarEntorno = (nuevo: string) => {
        setEntorno(nuevo)
        localStorage.setItem('entorno', nuevo)
        window.location.reload()
    }

    return (
        <div style={{ marginBottom: '1rem' }}>
            <label>🌐 Entorno actual: </label>
            <select value={entorno} onChange={(e) => cambiarEntorno(e.target.value)}>
                <option value="qa">QA</option>
                <option value="prod">PROD</option>
            </select>
        </div>
    )
}