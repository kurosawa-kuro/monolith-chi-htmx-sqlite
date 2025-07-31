import { defineConfig, devices } from '@playwright/test';

/**
 * Simple Playwright configuration for basic testing
 */
export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: false,
  retries: 0,
  workers: 1,
  reporter: 'html',
  use: {
    baseURL: 'http://localhost:8080',
    trace: 'off',
    screenshot: 'only-on-failure',
    video: 'off',
  },

  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],

  /* Run your local dev server before starting the tests */
  webServer: {
    command: 'make run-dev',
    url: 'http://localhost:8080',
    reuseExistingServer: true,
    timeout: 60 * 1000,
    stdout: 'pipe',
    stderr: 'pipe',
  },
}); 