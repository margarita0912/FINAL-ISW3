// tests/FacturaVisual.test.tsx
import { render, screen } from '@testing-library/react';
import React from 'react';

function FacturaVisual({ producto, precio }: { producto: string; precio: number }) {
    return (
        <div>
            <h2>Factura</h2>
            <p>Producto: {producto}</p>
            <p>Precio: ${precio}</p>
        </div>
    );
}

test('renderiza factura con datos correctos', () => {
    render(<FacturaVisual producto="Mouse" precio={100} />);
    expect(screen.getByText('Factura')).toBeInTheDocument();
    expect(screen.getByText('Producto: Mouse')).toBeInTheDocument();
    expect(screen.getByText('Precio: $100')).toBeInTheDocument();
});
