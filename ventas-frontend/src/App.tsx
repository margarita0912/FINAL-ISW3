import { Routes, Route, Link } from 'react-router-dom'
import EntornoSelector from './components/EntornoSelector'
import RoleRoute from './components/RoleRoute'

import Login from './pages/Login'
import Productos from './pages/Productos'
import Ventas from './pages/Ventas'
import CrearUsuario from './pages/CrearUsuario'
import CrearProductos from './pages/CrearProductos'

function App() {
    const rol = localStorage.getItem('rol') || '' // 'admin' | 'vendedor' | 'comprador' | 'precio' | ''

    return (
        <div style={{ padding: '2rem' }}>
            <h1>🛍️ Ventas App</h1>
            <EntornoSelector />
            <nav>
                <Link to="/">Productos</Link>
                {/* Ventas solo para vendedor */}
                {['vendedor'].includes(rol) && <> | <Link to="/ventas">Ventas</Link></>}
                <> | <Link to="/login">Login</Link></>

                {/* Crear Usuario: comprador o precio (según backend) */}
                {['comprador', 'precio'].includes(rol) && (
                    <> | <Link to="/crear-usuario">Crear Usuario</Link></>
                )}

                {/* Crear Producto: vendedor o comprador (según backend) */}
                {['vendedor', 'comprador'].includes(rol) && (
                    <> | <Link to="/crear-producto">Crear Producto</Link></>
                )}
            </nav>

            <Routes>
                <Route path="/" element={<Productos />} />

                {/* Ventas -> solo vendedor */}
                <Route
                    path="/ventas"
                    element={
                        <RoleRoute allowed={['vendedor']}>
                            <Ventas />
                        </RoleRoute>
                    }
                />

                <Route path="/login" element={<Login />} />

                {/* Crear Usuario -> comprador o precio */}
                <Route
                    path="/crear-usuario"
                    element={
                        <RoleRoute allowed={['comprador', 'precio']}>
                            <CrearUsuario />
                        </RoleRoute>
                    }
                />

                {/* Crear Producto -> vendedor o comprador */}
                <Route
                    path="/crear-producto"
                    element={
                        <RoleRoute allowed={['vendedor', 'comprador']}>
                            <CrearProductos />
                        </RoleRoute>
                    }
                />
            </Routes>
        </div>
    )
}

export default App
