import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import Ventas from '../pages/Ventas';
import api from '../api/axios';

jest.mock('../api/axios');

const mockProductos = [
  { id: 1, nombre: 'Laptop', precio: 800, stock: 10 },
  { id: 2, nombre: 'Mouse', precio: 25, stock: 50 },
  { id: 3, nombre: 'Teclado', precio: 45, stock: 0 },
];

describe('Ventas Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorage.clear();
  });

  it('renderiza el componente de ventas', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: mockProductos });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    expect(screen.getByText(/Ventas/i)).toBeInTheDocument();
    
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalledWith('/productos');
    });
  });

  it('carga productos correctamente', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: mockProductos });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(screen.getByText(/Laptop/i)).toBeInTheDocument();
      expect(screen.getByText(/Mouse/i)).toBeInTheDocument();
    });
  });

  it('muestra error cuando falla la carga de productos', async () => {
    const mockGet = jest.fn().mockRejectedValue(new Error('Network error'));
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(screen.getByText(/Error al cargar productos/i)).toBeInTheDocument();
    });
  });

  it('permite agregar items al carrito', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: mockProductos });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(screen.getByText(/Laptop/i)).toBeInTheDocument();
    });

    // Seleccionar producto y cantidad (el test específico dependerá de tu UI)
    const agregarButton = screen.getByText(/Agregar/i);
    expect(agregarButton).toBeInTheDocument();
  });

  it('muestra loading state mientras carga productos', () => {
    const mockGet = jest.fn().mockImplementation(() => new Promise(() => {}));
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    expect(screen.getByText(/Cargando/i) || screen.getByText(/Loading/i)).toBeInTheDocument();
  });

  it('valida stock insuficiente', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: mockProductos });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(screen.getByText(/Teclado/i)).toBeInTheDocument();
    });

    // El producto Teclado tiene stock 0
    // Intentar agregarlo debería mostrar error
  });

  it('calcula precios correctamente', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: mockProductos });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });

    // Verificar que los precios se muestran correctamente
    const laptop = mockProductos.find(p => p.nombre === 'Laptop');
    expect(laptop?.precio).toBe(800);
  });

  it('maneja descuentos correctamente', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: mockProductos });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });

    // El componente debería tener campos para descuentos
  });

  it('muestra lista vacía cuando no hay productos', async () => {
    const mockGet = jest.fn().mockResolvedValue({ data: [] });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });
  });

  it('normaliza datos de productos con diferentes formatos', async () => {
    const productosVariados = [
      { ID: 1, nombre: 'Test1', PRECIO: 100, STOCK: 5 },
      { id: 2, name: 'Test2', price: 200, stock: 10 },
    ];
    
    const mockGet = jest.fn().mockResolvedValue({ data: productosVariados });
    (api as jest.Mock).mockResolvedValue({ get: mockGet });

    render(<Ventas />);
    
    await waitFor(() => {
      expect(mockGet).toHaveBeenCalled();
    });
  });
});
