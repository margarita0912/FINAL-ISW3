export default {
    preset: 'ts-jest',
    testEnvironment: 'jest-environment-jsdom',
    setupFiles: ['./jest.setup.js'],
    setupFilesAfterEnv: ['./setupTests.ts'],
    moduleNameMapper: {
        '^@/(.*)$': '<rootDir>/sr/$1',
        // Mock MSW completamente para evitar problemas ESM
        '^msw/node$': '<rootDir>/__mocks__/msw.js',
        // Mock api/axios para evitar problemas con import.meta
        '^../api/axios$': '<rootDir>/src/api/__mocks__/axios.ts',
        '^@/api/axios$': '<rootDir>/src/api/__mocks__/axios.ts',
    },
    transform: {
        '^.+\\.tsx?$': ['ts-jest', {
            tsconfig: {
                module: 'esnext',
                target: 'esnext',
            },
        }],
    },
    transformIgnorePatterns: ["node_modules/"],
    modulePathIgnorePatterns: ['/dist/'],
    // Forzar a Jest a usar .tsx en lugar de .js
    moduleFileExtensions: ['ts', 'tsx', 'js', 'jsx', 'json'],
    collectCoverage: true,
    coverageDirectory: 'coverage',
    testMatch: ['**/tests/**/*.test.ts?(x)', '**/src/tests/**/*.test.ts?(x)'],
    globals: {
        'import.meta': {
            env: {
                VITE_API_URL: 'http://localhost:8080',
                MODE: 'test',
            },
        },
    },
};
