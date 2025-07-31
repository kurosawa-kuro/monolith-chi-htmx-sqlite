import { test, expect } from '@playwright/test';

test.describe('Category Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should create a new category', async ({ page }) => {
    const categoryName = `Test Category ${Date.now()}`;
    
    // Fill in the category form
    await page.locator('[data-testid="category-name-input"]').fill(categoryName);
    
    // Submit the form
    await page.locator('[data-testid="create-category-button"]').click();
    
    // Wait for the category to be created
    await page.waitForResponse(response => 
      response.url().includes('/categories') && response.status() === 303
    );
    
    // Verify the category appears in the filter dropdown
    await expect(page.locator(`text=${categoryName}`)).toBeVisible();
  });

  test('should delete a category', async ({ page }) => {
    // Find the first category in the list
    const firstCategory = page.locator('[data-testid="category-item"]').first();
    const categoryName = await firstCategory.locator('[data-testid="category-name"]').textContent();
    
    // Click delete button
    await firstCategory.locator('[data-testid="delete-category-button"]').click();
    
    // Wait for the delete to complete
    await page.waitForResponse(response => 
      response.url().includes('/categories') && response.method() === 'DELETE'
    );
    
    // Verify the category is removed
    await expect(page.locator(`text=${categoryName}`)).not.toBeVisible();
  });

  test('should display category count', async ({ page }) => {
    // Check if category count is displayed
    const categoryItems = page.locator('[data-testid="category-item"]');
    const count = await categoryItems.count();
    
    // Verify category count is shown
    await expect(page.locator('[data-testid="category-count"]')).toContainText(count.toString());
  });

  test('should filter todos by category', async ({ page }) => {
    // Get all available categories
    const categoryOptions = page.locator('[data-testid="category-filter"] option');
    const categoryCount = await categoryOptions.count();
    
    if (categoryCount > 1) {
      // Select a specific category
      const selectedCategory = await categoryOptions.nth(1).textContent();
      await page.locator('[data-testid="category-filter"]').selectOption({ index: 1 });
      
      // Wait for the filter to be applied
      await page.waitForResponse(response => 
        response.url().includes('category=') && response.status() === 200
      );
      
      // Verify the filter is applied
      await expect(page).toHaveURL(/category=/);
      
      // Verify only todos from the selected category are shown
      const todoItems = page.locator('[data-testid="todo-item"]');
      for (let i = 0; i < await todoItems.count(); i++) {
        const todoCategories = todoItems.nth(i).locator('[data-testid="todo-categories"]');
        await expect(todoCategories).toContainText(selectedCategory!);
      }
    }
  });
}); 