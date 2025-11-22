import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';

function FormularioVenta() {
    const [producto, setProducto] = React.useState('');
    const [precio, setPrecio] = React.useState(0);

    return (
        <form>
            <label>
                Producto:
                <input
                    value={producto}
                    onChange={(e) => setProducto(e.target.value)}
                    placeholder="Nombre del producto"
                />
            </label>
            <label>
                Precio:
                <input
                    type="number"
                    value={precio}
                    onChange={(e) => setPrecio(Number(e.target.value))}
                    placeholder="Precio"
                />
            </label>
            <button type="submit">Enviar</button>
        </form>
    );
}

test('permite ingresar producto y precio', () => {
    render(<FormularioVenta />);
    const inputProducto = screen.getByPlaceholderText('Nombre del producto');
    const inputPrecio = screen.getByPlaceholderText('Precio');

    fireEvent.change(inputProducto, { target: { value: 'Teclado' } });
    fireEvent.change(inputPrecio, { target: { value: '2500' } });

    expect(inputProducto).toHaveValue('Teclado');
    expect(inputPrecio).toHaveValue(2500);
});
