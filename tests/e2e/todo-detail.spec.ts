import { test, expect } from '@playwright/test';

test.describe('Todo Detail Page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should display todo detail page', async ({ page }) => {
    // Create a test todo first
    const todoTitle = `Test Todo Detail ${Date.now()}`;
    await page.locator('#title').fill(todoTitle);
    await page.locator('button[type="submit"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Click on the todo to go to detail page
    await page.locator(`text=${todoTitle}`).click();
    
    // Verify we're on the detail page
    await expect(page).toHaveURL(/\/todos\/\d+/);
    await expect(page.locator('h1:has-text("Todo詳細")')).toBeVisible();
    await expect(page.locator(`text=${todoTitle}`)).toBeVisible();
  });

  test('should display todo information correctly', async ({ page }) => {
    // Create a test todo first
    const todoTitle = `Test Todo Info ${Date.now()}`;
    await page.locator('#title').fill(todoTitle);
    await page.locator('button[type="submit"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Click on the todo to go to detail page
    await page.locator(`text=${todoTitle}`).click();
    
    // Verify todo information is displayed
    await expect(page.locator('text=タイトル')).toBeVisible();
    await expect(page.locator('text=ステータス')).toBeVisible();
    await expect(page.locator(`text=${todoTitle}`)).toBeVisible();
    await expect(page.locator('text=未完了')).toBeVisible();
  });

  test('should display categories section', async ({ page }) => {
    // Create a test todo first
    const todoTitle = `Test Todo Categories ${Date.now()}`;
    await page.locator('#title').fill(todoTitle);
    await page.locator('button[type="submit"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Click on the todo to go to detail page
    await page.locator(`text=${todoTitle}`).click();
    
    // Verify categories section is displayed
    await expect(page.locator('text=カテゴリー')).toBeVisible();
  });

  test('should have edit functionality', async ({ page }) => {
    // Create a test todo first
    const todoTitle = `Test Todo Edit ${Date.now()}`;
    await page.locator('#title').fill(todoTitle);
    await page.locator('button[type="submit"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Click on the todo to go to detail page
    await page.locator(`text=${todoTitle}`).click();
    
    // Verify edit button is present
    await expect(page.locator('button:has-text("編集")')).toBeVisible();
  });

  test('should have back button', async ({ page }) => {
    // Create a test todo first
    const todoTitle = `Test Todo Back ${Date.now()}`;
    await page.locator('#title').fill(todoTitle);
    await page.locator('button[type="submit"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Click on the todo to go to detail page
    await page.locator(`text=${todoTitle}`).click();
    
    // Verify back button is present
    await expect(page.locator('a:has-text("戻る")')).toBeVisible();
    
    // Click back button and verify we return to home page
    await page.locator('a:has-text("戻る")').click();
    await expect(page).toHaveURL('/');
  });

  test('should handle todo not found', async ({ page }) => {
    // Try to access a non-existent todo
    await page.goto('/todos/999999');
    
    // Should show 404 or error page
    await expect(page.locator('text=404')).toBeVisible();
  });

  test('should display created and updated timestamps', async ({ page }) => {
    // Create a test todo first
    const todoTitle = `Test Todo Timestamps ${Date.now()}`;
    await page.locator('#title').fill(todoTitle);
    await page.locator('button[type="submit"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Click on the todo to go to detail page
    await page.locator(`text=${todoTitle}`).click();
    
    // Verify timestamps are displayed
    await expect(page.locator('text=作成日時')).toBeVisible();
    await expect(page.locator('text=更新日時')).toBeVisible();
  });
}); 