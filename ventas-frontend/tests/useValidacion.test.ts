// tests/useValidacion.test.ts
global.fetch = jest.fn(() =>
    Promise.resolve({
        ok: true,
        status: 200,
        json: async () => ({ valid: true }),
    } as unknown as Response)
);

async function validarToken(token: string): Promise<boolean> {
    const res = await fetch(`/validar?token=${token}`);
    const data = await res.json();
    return data.valid;
}

test('valida token correctamente', async () => {
    const resultado = await validarToken('abc123');
    expect(resultado).toBe(true);
});
