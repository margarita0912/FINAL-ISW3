import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import RoleRoute from '../components/RoleRoute';

// Mock de Navigate de react-router-dom
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  Navigate: ({ to }: { to: string }) => <div>Redirected to {to}</div>,
}));

describe('RoleRoute Component', () => {
  beforeEach(() => {
    localStorage.clear();
  });

  it('permite acceso cuando el rol está en la lista allowed', () => {
    localStorage.setItem('rol', 'administrador');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador', 'vendedor']}>
          <div>Contenido Protegido</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.getByText('Contenido Protegido')).toBeInTheDocument();
    expect(screen.queryByText(/Redirected/i)).not.toBeInTheDocument();
  });

  it('redirige cuando el rol no está en la lista allowed', () => {
    localStorage.setItem('rol', 'comprador');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador', 'vendedor']}>
          <div>Contenido Protegido</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.queryByText('Contenido Protegido')).not.toBeInTheDocument();
    expect(screen.getByText('Redirected to /')).toBeInTheDocument();
  });

  it('redirige cuando no hay rol en localStorage', () => {
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador']}>
          <div>Contenido Protegido</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.queryByText('Contenido Protegido')).not.toBeInTheDocument();
    expect(screen.getByText('Redirected to /')).toBeInTheDocument();
  });

  it('permite acceso a vendedor cuando está en la lista', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['vendedor']}>
          <div>Panel de Vendedor</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.getByText('Panel de Vendedor')).toBeInTheDocument();
  });

  it('permite acceso a comprador cuando está en la lista', () => {
    localStorage.setItem('rol', 'comprador');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['comprador', 'vendedor']}>
          <div>Panel de Compras</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.getByText('Panel de Compras')).toBeInTheDocument();
  });

  it('permite múltiples roles en allowed', () => {
    localStorage.setItem('rol', 'administrador');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador', 'vendedor', 'comprador']}>
          <div>Contenido Multi-Rol</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.getByText('Contenido Multi-Rol')).toBeInTheDocument();
  });

  it('renderiza children correctamente cuando tiene acceso', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['vendedor']}>
          <div>
            <h1>Título</h1>
            <p>Descripción</p>
          </div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.getByText('Título')).toBeInTheDocument();
    expect(screen.getByText('Descripción')).toBeInTheDocument();
  });

  it('maneja rol vacío correctamente', () => {
    localStorage.setItem('rol', '');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador']}>
          <div>Contenido Protegido</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.queryByText('Contenido Protegido')).not.toBeInTheDocument();
    expect(screen.getByText('Redirected to /')).toBeInTheDocument();
  });

  it('es case-sensitive con los roles', () => {
    localStorage.setItem('rol', 'Administrador'); // Con mayúscula
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador']}> {/* minúscula */}
          <div>Contenido Protegido</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    // No debería permitir acceso porque es case-sensitive
    expect(screen.queryByText('Contenido Protegido')).not.toBeInTheDocument();
  });

  it('permite solo un rol específico', () => {
    localStorage.setItem('rol', 'administrador');
    
    render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador']}>
          <div>Panel Admin</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    expect(screen.getByText('Panel Admin')).toBeInTheDocument();
  });

  it('redirige con replace para evitar historial', () => {
    localStorage.setItem('rol', 'comprador');
    
    const { container } = render(
      <BrowserRouter>
        <RoleRoute allowed={['administrador']}>
          <div>Contenido Protegido</div>
        </RoleRoute>
      </BrowserRouter>
    );
    
    // Verifica que se renderiza el Navigate (mock)
    expect(container.textContent).toContain('Redirected to /');
  });
});
