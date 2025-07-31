import { test, expect } from '@playwright/test';

test.describe('Todo CRUD Operations', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/');
  });

  test('should create a new todo', async ({ page }) => {
    const todoTitle = `Test Todo ${Date.now()}`;
    
    // Fill in the todo form
    await page.locator('[data-testid="todo-title-input"]').fill(todoTitle);
    
    // Submit the form
    await page.locator('[data-testid="create-todo-button"]').click();
    
    // Wait for the todo to be created
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Verify the todo appears in the list
    await expect(page.locator(`text=${todoTitle}`)).toBeVisible();
  });

  test('should update todo status', async ({ page }) => {
    // Find the first todo and click its status toggle
    const firstTodo = page.locator('[data-testid="todo-item"]').first();
    const statusToggle = firstTodo.locator('[data-testid="status-toggle"]');
    
    // Get initial status
    const initialStatus = await statusToggle.getAttribute('data-status');
    
    // Toggle status
    await statusToggle.click();
    
    // Wait for the update to complete
    await page.waitForResponse(response => 
      response.url().includes('/status') && response.status() === 200
    );
    
    // Verify status changed
    const newStatus = await statusToggle.getAttribute('data-status');
    expect(newStatus).not.toBe(initialStatus);
  });

  test('should delete a todo', async ({ page }) => {
    // Find the first todo
    const firstTodo = page.locator('[data-testid="todo-item"]').first();
    const todoTitle = await firstTodo.locator('[data-testid="todo-title"]').textContent();
    
    // Click delete button
    await firstTodo.locator('[data-testid="delete-todo-button"]').click();
    
    // Wait for the delete to complete
    await page.waitForResponse(response => 
      response.url().includes('/todos') && response.method() === 'DELETE'
    );
    
    // Verify the todo is removed
    await expect(page.locator(`text=${todoTitle}`)).not.toBeVisible();
  });

  test('should filter todos by category', async ({ page }) => {
    // Select a category from the filter
    const categorySelect = page.locator('[data-testid="category-filter"]');
    await categorySelect.selectOption({ index: 1 }); // Select first category
    
    // Wait for the filter to be applied
    await page.waitForResponse(response => 
      response.url().includes('category=') && response.status() === 200
    );
    
    // Verify the URL contains the category filter
    await expect(page).toHaveURL(/category=/);
  });

  test('should handle pagination', async ({ page }) => {
    // Check if pagination controls exist
    const pagination = page.locator('[data-testid="pagination"]');
    await expect(pagination).toBeVisible();
    
    // Click next page if available
    const nextButton = pagination.locator('[data-testid="next-page"]');
    if (await nextButton.isVisible()) {
      await nextButton.click();
      
      // Wait for the page to load
      await page.waitForResponse(response => 
        response.url().includes('page=') && response.status() === 200
      );
      
      // Verify URL contains page parameter
      await expect(page).toHaveURL(/page=/);
    }
  });
}); 