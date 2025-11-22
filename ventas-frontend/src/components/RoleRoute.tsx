// src/components/RoleRoute.tsx
import { Navigate } from 'react-router-dom'
import type { ReactNode } from 'react'  // 👈 usa "type" al importar solo tipos

export default function RoleRoute({
                                      allowed,
                                      children,
                                  }: {
    allowed: string[]
    children: ReactNode
}) {
    const rol: string = localStorage.getItem('rol') || ''
    return allowed.includes(rol) ? <>{children}</> : <Navigate to="/" replace />
}
