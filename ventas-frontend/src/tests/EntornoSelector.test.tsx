import { render, screen } from '@testing-library/react';
import EntornoSelector from '../components/EntornoSelector';

describe('EntornoSelector Component', () => {
  it('renderiza sin errores', () => {
    const { container } = render(<EntornoSelector />);
    expect(container).toBeInTheDocument();
  });

  it('retorna null como se espera', () => {
    const { container } = render(<EntornoSelector />);
    expect(container.firstChild).toBeNull();
  });

  it('no muestra ningún elemento visual', () => {
    render(<EntornoSelector />);
    // No debería haber texto o elementos visibles
    expect(screen.queryByRole('button')).not.toBeInTheDocument();
    expect(screen.queryByRole('combobox')).not.toBeInTheDocument();
  });

  it('mantiene compatibilidad con imports existentes', () => {
    // Este test verifica que el componente puede ser importado sin errores
    expect(EntornoSelector).toBeDefined();
    expect(typeof EntornoSelector).toBe('function');
  });

  it('puede ser renderizado múltiples veces sin errores', () => {
    const { rerender } = render(<EntornoSelector />);
    expect(() => rerender(<EntornoSelector />)).not.toThrow();
  });
});
