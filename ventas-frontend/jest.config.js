export default {
    preset: 'ts-jest',
    testEnvironment: 'jest-environment-jsdom',
    setupFiles: ['./jest.setup.js'], // ← nuevo
    setupFilesAfterEnv: ['./setupTests.ts'],
    moduleNameMapper: {
        '^@/(.*)$': '<rootDir>/sr/$1',
    },
    transform: {
        '^.+\\.tsx?$': 'ts-jest',
    },
    // Transformar msw y sus dependencias ESM que vienen en node_modules
    transformIgnorePatterns: ["node_modules/(?!(msw|@mswjs|until-async)/)"],
    collectCoverage: true,
    coverageDirectory: 'coverage',
    testMatch: ['**/tests/**/*.test.ts?(x)'],
};
