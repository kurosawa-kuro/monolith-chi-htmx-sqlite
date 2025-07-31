import { test, expect } from '@playwright/test';

test.describe('Home Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display the home page', async ({ page }) => {
    // Check page title
    await expect(page).toHaveTitle(/Todo App/);
    
    // Check main heading
    await expect(page.locator('h1')).toContainText('Todo');
    
    // Check if the page loads without errors
    await expect(page.locator('body')).toBeVisible();
  });

  test('should display todo list', async ({ page }) => {
    // Check if todo list container exists
    await expect(page.locator('[data-testid="todo-list"]')).toBeVisible();
    
    // Check if pagination exists
    await expect(page.locator('[data-testid="pagination"]')).toBeVisible();
  });

  test('should display category filter', async ({ page }) => {
    // Check if category filter exists
    await expect(page.locator('[data-testid="category-filter"]')).toBeVisible();
  });

  test('should have working navigation', async ({ page }) => {
    // Check if all main navigation elements are present
    await expect(page.locator('nav')).toBeVisible();
  });
}); 