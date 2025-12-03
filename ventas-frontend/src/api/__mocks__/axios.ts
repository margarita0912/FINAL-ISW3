// Mock de api/axios para Jest - evita problemas con import.meta
const mockAxiosInstance = {
  get: jest.fn(() => Promise.resolve({ data: [] })),
  post: jest.fn(() => Promise.resolve({ data: {} })),
  put: jest.fn(() => Promise.resolve({ data: {} })),
  delete: jest.fn(() => Promise.resolve({ data: {} })),
};

export default jest.fn(() => Promise.resolve(mockAxiosInstance));
