import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import CrearUsuario from '../pages/CrearUsuario';
import api from '../api/axios';

jest.mock('../api/axios');

describe('CrearUsuario Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorage.clear();
  });

  it('muestra mensaje de sin permisos para usuarios no autorizados', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearUsuario />);
    
    expect(screen.getByText(/No tenés permisos/i)).toBeInTheDocument();
  });

  it('permite acceso a usuarios con rol "precio"', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    expect(screen.getByRole('heading', { name: /Crear Usuario/i })).toBeInTheDocument();
    expect(screen.queryByText(/No tenés permisos/i)).not.toBeInTheDocument();
  });

  it('renderiza todos los campos del formulario', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    expect(screen.getByPlaceholderText(/lucas/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/clave123/i)).toBeInTheDocument();
    expect(screen.getByRole('combobox')).toBeInTheDocument();
  });

  it('actualiza el input de nombre correctamente', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    const nombreInput = screen.getByPlaceholderText(/lucas/i) as HTMLInputElement;
    fireEvent.change(nombreInput, { target: { value: 'juan' } });
    
    expect(nombreInput.value).toBe('juan');
  });

  it('actualiza el input de clave correctamente', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    const claveInput = screen.getByPlaceholderText(/clave123/i) as HTMLInputElement;
    fireEvent.change(claveInput, { target: { value: 'password456' } });
    
    expect(claveInput.value).toBe('password456');
  });

  it('permite cambiar el rol seleccionado', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    const rolSelect = screen.getByRole('combobox') as HTMLSelectElement;
    fireEvent.change(rolSelect, { target: { value: 'vendedor' } });
    
    expect(rolSelect.value).toBe('vendedor');
  });

  it('muestra todas las opciones de rol', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    expect(screen.getByText(/Comprador/i)).toBeInTheDocument();
    expect(screen.getByText(/Vendedor/i)).toBeInTheDocument();
    expect(screen.getByText(/Precio/i)).toBeInTheDocument();
  });

  it('crea usuario exitosamente', async () => {
    localStorage.setItem('rol', 'precio');
    localStorage.setItem('token', 'fake-token');
    
    const mockPost = jest.fn().mockResolvedValue({ data: { id: 1 } });
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<CrearUsuario />);
    
    const nombreInput = screen.getByPlaceholderText(/lucas/i);
    const claveInput = screen.getByPlaceholderText(/clave123/i);
    const crearButton = screen.getByRole('button', { name: /Crear Usuario/i });

    fireEvent.change(nombreInput, { target: { value: 'testuser' } });
    fireEvent.change(claveInput, { target: { value: 'testpass' } });
    fireEvent.click(crearButton);

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith(
        '/usuarios',
        { nombre: 'testuser', clave: 'testpass', rol: 'comprador' },
        { headers: { Authorization: 'Bearer fake-token' } }
      );
    });

    expect(screen.getByText(/Usuario creado con éxito/i)).toBeInTheDocument();
  });

  it('muestra error cuando falla la creación', async () => {
    localStorage.setItem('rol', 'precio');
    localStorage.setItem('token', 'fake-token');
    
    const mockPost = jest.fn().mockRejectedValue(new Error('Server error'));
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<CrearUsuario />);
    
    const nombreInput = screen.getByPlaceholderText(/lucas/i);
    const claveInput = screen.getByPlaceholderText(/clave123/i);
    const crearButton = screen.getByRole('button', { name: /Crear Usuario/i });

    fireEvent.change(nombreInput, { target: { value: 'test' } });
    fireEvent.change(claveInput, { target: { value: 'pass' } });
    fireEvent.click(crearButton);

    await waitFor(() => {
      expect(screen.getByText(/Error al crear usuario/i)).toBeInTheDocument();
    });
  });

  it('limpia el formulario después de crear exitosamente', async () => {
    localStorage.setItem('rol', 'precio');
    localStorage.setItem('token', 'fake-token');
    
    const mockPost = jest.fn().mockResolvedValue({ data: { id: 1 } });
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<CrearUsuario />);
    
    const nombreInput = screen.getByPlaceholderText(/lucas/i) as HTMLInputElement;
    const claveInput = screen.getByPlaceholderText(/clave123/i) as HTMLInputElement;
    const crearButton = screen.getByRole('button', { name: /Crear Usuario/i });

    fireEvent.change(nombreInput, { target: { value: 'test' } });
    fireEvent.change(claveInput, { target: { value: 'pass' } });
    fireEvent.click(crearButton);

    await waitFor(() => {
      expect(nombreInput.value).toBe('');
      expect(claveInput.value).toBe('');
    });
  });

  it('el campo de clave es de tipo password', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    const claveInput = screen.getByPlaceholderText(/clave123/i);
    expect(claveInput).toHaveAttribute('type', 'password');
  });

  it('rol por defecto es comprador', () => {
    localStorage.setItem('rol', 'precio');
    
    render(<CrearUsuario />);
    
    const rolSelect = screen.getByRole('combobox') as HTMLSelectElement;
    expect(rolSelect.value).toBe('comprador');
  });
});
