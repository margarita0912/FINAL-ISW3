import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import CrearProductos from '../pages/CrearProductos';
import api from '../api/axios';

jest.mock('../api/axios');

describe('CrearProductos Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorage.clear();
  });

  it('muestra mensaje de sin permisos para usuarios no autorizados', () => {
    localStorage.setItem('rol', 'administrador');
    
    render(<CrearProductos />);
    
    expect(screen.getByText(/No tenés permisos/i)).toBeInTheDocument();
  });

  it('permite acceso a vendedores', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    expect(screen.getByRole('heading', { name: /Crear Producto/i })).toBeInTheDocument();
    expect(screen.queryByText(/No tenés permisos/i)).not.toBeInTheDocument();
  });

  it('permite acceso a compradores', () => {
    localStorage.setItem('rol', 'comprador');
    
    render(<CrearProductos />);
    
    expect(screen.getByRole('heading', { name: /Crear Producto/i })).toBeInTheDocument();
  });

  it('renderiza todos los campos del formulario', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    expect(screen.getByPlaceholderText(/Camisa/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/2500/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/10/i)).toBeInTheDocument();
  });

  it('actualiza el input de nombre correctamente', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    const nombreInput = screen.getByPlaceholderText(/Camisa/i) as HTMLInputElement;
    fireEvent.change(nombreInput, { target: { value: 'Laptop' } });
    
    expect(nombreInput.value).toBe('Laptop');
  });

  it('actualiza el input de precio correctamente', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    const precioInput = screen.getByPlaceholderText(/2500/i) as HTMLInputElement;
    fireEvent.change(precioInput, { target: { value: '1500' } });
    
    expect(precioInput.value).toBe('1500');
  });

  it('actualiza el input de stock correctamente', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    const stockInput = screen.getByPlaceholderText(/10/i) as HTMLInputElement;
    fireEvent.change(stockInput, { target: { value: '25' } });
    
    expect(stockInput.value).toBe('25');
  });

  it('muestra error cuando faltan campos', async () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    const crearButton = screen.getByRole('button', { name: /Crear Producto/i });
    fireEvent.click(crearButton);
    
    await waitFor(() => {
      expect(screen.getByText(/Todos los campos son obligatorios/i)).toBeInTheDocument();
    });
  });

  it('crea producto exitosamente', async () => {
    localStorage.setItem('rol', 'vendedor');
    localStorage.setItem('token', 'fake-token');
    
    const mockPost = jest.fn().mockResolvedValue({ data: { id: 1 } });
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<CrearProductos />);
    
    const nombreInput = screen.getByPlaceholderText(/Camisa/i);
    const precioInput = screen.getByPlaceholderText(/2500/i);
    const stockInput = screen.getByPlaceholderText(/10/i);
    const crearButton = screen.getByRole('button', { name: /Crear Producto/i });

    fireEvent.change(nombreInput, { target: { value: 'Mouse' } });
    fireEvent.change(precioInput, { target: { value: '25' } });
    fireEvent.change(stockInput, { target: { value: '50' } });
    fireEvent.click(crearButton);

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith(
        '/productos',
        { nombre: 'Mouse', precio: 25, stock: 50 },
        { headers: { Authorization: 'Bearer fake-token' } }
      );
    });

    expect(screen.getByText(/Producto creado con éxito/i)).toBeInTheDocument();
  });

  it('muestra error cuando falla la creación', async () => {
    localStorage.setItem('rol', 'vendedor');
    localStorage.setItem('token', 'fake-token');
    
    const mockPost = jest.fn().mockRejectedValue(new Error('Server error'));
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<CrearProductos />);
    
    const nombreInput = screen.getByPlaceholderText(/Camisa/i);
    const precioInput = screen.getByPlaceholderText(/2500/i);
    const stockInput = screen.getByPlaceholderText(/10/i);
    const crearButton = screen.getByRole('button', { name: /Crear Producto/i });

    fireEvent.change(nombreInput, { target: { value: 'Test' } });
    fireEvent.change(precioInput, { target: { value: '100' } });
    fireEvent.change(stockInput, { target: { value: '10' } });
    fireEvent.click(crearButton);

    await waitFor(() => {
      expect(screen.getByText(/Error al crear producto/i)).toBeInTheDocument();
    });
  });

  it('limpia el formulario después de crear exitosamente', async () => {
    localStorage.setItem('rol', 'vendedor');
    localStorage.setItem('token', 'fake-token');
    
    const mockPost = jest.fn().mockResolvedValue({ data: { id: 1 } });
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<CrearProductos />);
    
    const nombreInput = screen.getByPlaceholderText(/Camisa/i) as HTMLInputElement;
    const precioInput = screen.getByPlaceholderText(/2500/i) as HTMLInputElement;
    const stockInput = screen.getByPlaceholderText(/10/i) as HTMLInputElement;
    const crearButton = screen.getByRole('button', { name: /Crear Producto/i });

    fireEvent.change(nombreInput, { target: { value: 'Test' } });
    fireEvent.change(precioInput, { target: { value: '100' } });
    fireEvent.change(stockInput, { target: { value: '10' } });
    fireEvent.click(crearButton);

    await waitFor(() => {
      expect(nombreInput.value).toBe('');
      expect(precioInput.value).toBe('');
      expect(stockInput.value).toBe('');
    });
  });

  it('valida números negativos en precio', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    const precioInput = screen.getByPlaceholderText(/2500/i) as HTMLInputElement;
    expect(precioInput).toHaveAttribute('min', '0');
  });

  it('valida números negativos en stock', () => {
    localStorage.setItem('rol', 'vendedor');
    
    render(<CrearProductos />);
    
    const stockInput = screen.getByPlaceholderText(/10/i) as HTMLInputElement;
    expect(stockInput).toHaveAttribute('min', '0');
  });
});
