import { test, expect } from '@playwright/test';

test.describe('Basic Application Test', () => {
  test('should load the application', async ({ page }) => {
    // Navigate to the application
    await page.goto('/');
    
    // Check if the page loads
    await expect(page).toHaveTitle(/Todo/);
    
    // Check if the body is visible
    await expect(page.locator('body')).toBeVisible();
    
    // Take a screenshot for verification
    await page.screenshot({ path: 'test-results/screenshots/basic-load.png' });
  });

  test('should have basic HTML structure', async ({ page }) => {
    await page.goto('/');
    
    // Check for basic HTML elements
    await expect(page.locator('html')).toBeVisible();
    await expect(page.locator('body')).toBeVisible();
    
    // Check if there's some content
    const bodyText = await page.locator('body').textContent();
    expect(bodyText).toBeTruthy();
    expect(bodyText!.length).toBeGreaterThan(0);
    
    // Check if title is present
    await expect(page).toHaveTitle(/Todo/);
  });
}); 