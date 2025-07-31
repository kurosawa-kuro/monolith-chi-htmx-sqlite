import { test, expect } from '@playwright/test';

test('minimal test - should load page', async ({ page }) => {
  await page.goto('/');
  await expect(page.locator('body')).toBeVisible();
  console.log('✅ Page loaded successfully');
}); 