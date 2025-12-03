import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import '@testing-library/jest-dom';
import Login from '../pages/Login';
import api from '../api/axios';

// Mock del módulo api
jest.mock('../api/axios');

describe('Login Component', () => {
  beforeEach(() => {
    // Limpiar mocks y localStorage
    jest.clearAllMocks();
    localStorage.clear();
    delete (window as any).location;
    (window as any).location = { href: '' };
  });

  it('renderiza el formulario de login', () => {
    render(<Login />);
    
    expect(screen.getByText(/Login/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/Email/i)).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/Password/i)).toBeInTheDocument();
    expect(screen.getByText(/Ingresar/i)).toBeInTheDocument();
  });

  it('actualiza los inputs correctamente', () => {
    render(<Login />);
    
    const emailInput = screen.getByPlaceholderText(/Email/i) as HTMLInputElement;
    const passwordInput = screen.getByPlaceholderText(/Password/i) as HTMLInputElement;

    fireEvent.change(emailInput, { target: { value: 'test@test.com' } });
    fireEvent.change(passwordInput, { target: { value: 'password123' } });

    expect(emailInput.value).toBe('test@test.com');
    expect(passwordInput.value).toBe('password123');
  });

  it('login exitoso guarda token y redirige', async () => {
    const mockPost = jest.fn().mockResolvedValue({
      data: { token: 'fake-token-123', rol: 'admin' }
    });
    
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<Login />);
    
    const emailInput = screen.getByPlaceholderText(/Email/i);
    const passwordInput = screen.getByPlaceholderText(/Password/i);
    const loginButton = screen.getByText(/Ingresar/i);

    fireEvent.change(emailInput, { target: { value: 'admin@test.com' } });
    fireEvent.change(passwordInput, { target: { value: 'admin123' } });
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(mockPost).toHaveBeenCalledWith('/login', { 
        nombre: 'admin@test.com', 
        clave: 'admin123' 
      });
    });

    expect(localStorage.getItem('token')).toBe('fake-token-123');
    expect(localStorage.getItem('rol')).toBe('admin');
    expect(window.location.href).toContain('/');
  });

  it('muestra error cuando el login falla', async () => {
    const mockPost = jest.fn().mockRejectedValue(new Error('Invalid credentials'));
    (api as jest.Mock).mockResolvedValue({ post: mockPost });

    render(<Login />);
    
    const emailInput = screen.getByPlaceholderText(/Email/i);
    const passwordInput = screen.getByPlaceholderText(/Password/i);
    const loginButton = screen.getByText(/Ingresar/i);

    fireEvent.change(emailInput, { target: { value: 'wrong@test.com' } });
    fireEvent.change(passwordInput, { target: { value: 'wrongpass' } });
    fireEvent.click(loginButton);

    await waitFor(() => {
      expect(screen.getByText(/Credenciales inválidas/i)).toBeInTheDocument();
    });

    expect(localStorage.getItem('token')).toBeNull();
  });

  it('no muestra error inicialmente', () => {
    render(<Login />);
    
    expect(screen.queryByText(/Credenciales inválidas/i)).not.toBeInTheDocument();
  });
});
