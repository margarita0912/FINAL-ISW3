import { render, screen } from '@testing-library/react';
import '@testing-library/jest-dom';
import { BrowserRouter } from 'react-router-dom';
import App from '../App';

describe('App Component', () => {
  it('renderiza sin crashes', () => {
    render(
      <BrowserRouter>
        <App />
      </BrowserRouter>
    );
    expect(screen.getByText(/Ventas App/i)).toBeInTheDocument();
  });

  it('muestra el título de la aplicación', () => {
    render(
      <BrowserRouter>
        <App />
      </BrowserRouter>
    );
    expect(screen.getByText(/Ventas App/i)).toBeInTheDocument();
  });

  it('muestra navegación', () => {
    render(
      <BrowserRouter>
        <App />
      </BrowserRouter>
    );
    const productosLinks = screen.getAllByText(/Productos/i);
    const loginLinks = screen.getAllByText(/Login/i);
    expect(productosLinks.length).toBeGreaterThan(0);
    expect(loginLinks.length).toBeGreaterThan(0);
  });
});
