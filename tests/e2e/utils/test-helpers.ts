import { Page, expect } from '@playwright/test';

export class TestHelpers {
  constructor(private page: Page) {}

  /**
   * Wait for the application to be ready
   */
  async waitForAppReady() {
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForSelector('body', { state: 'visible' });
  }

  /**
   * Create a test todo
   */
  async createTodo(title: string, categoryIds: number[] = []) {
    await this.page.locator('[data-testid="todo-title-input"]').fill(title);
    
    // Select categories if provided
    for (const categoryId of categoryIds) {
      await this.page.locator(`[data-testid="category-checkbox-${categoryId}"]`).check();
    }
    
    await this.page.locator('[data-testid="create-todo-button"]').click();
    
    // Wait for the todo to be created
    await this.page.waitForResponse(response => 
      response.url().includes('/todos') && response.status() === 303
    );
    
    // Verify the todo appears
    await expect(this.page.locator(`text=${title}`)).toBeVisible();
  }

  /**
   * Create a test category
   */
  async createCategory(name: string) {
    await this.page.locator('[data-testid="category-name-input"]').fill(name);
    await this.page.locator('[data-testid="create-category-button"]').click();
    
    // Wait for the category to be created
    await this.page.waitForResponse(response => 
      response.url().includes('/categories') && response.status() === 303
    );
    
    // Verify the category appears
    await expect(this.page.locator(`text=${name}`)).toBeVisible();
  }

  /**
   * Delete a todo by title
   */
  async deleteTodo(title: string) {
    const todoItem = this.page.locator(`[data-testid="todo-item"]:has-text("${title}")`);
    await todoItem.locator('[data-testid="delete-todo-button"]').click();
    
    // Wait for the delete to complete
    await this.page.waitForResponse(response => 
      response.url().includes('/todos') && response.method() === 'DELETE'
    );
    
    // Verify the todo is removed
    await expect(this.page.locator(`text=${title}`)).not.toBeVisible();
  }

  /**
   * Delete a category by name
   */
  async deleteCategory(name: string) {
    const categoryItem = this.page.locator(`[data-testid="category-item"]:has-text("${name}")`);
    await categoryItem.locator('[data-testid="delete-category-button"]').click();
    
    // Wait for the delete to complete
    await this.page.waitForResponse(response => 
      response.url().includes('/categories') && response.method() === 'DELETE'
    );
    
    // Verify the category is removed
    await expect(this.page.locator(`text=${name}`)).not.toBeVisible();
  }

  /**
   * Filter todos by category
   */
  async filterByCategory(categoryName: string) {
    await this.page.locator('[data-testid="category-filter"]').selectOption({ label: categoryName });
    
    // Wait for the filter to be applied
    await this.page.waitForResponse(response => 
      response.url().includes('category=') && response.status() === 200
    );
    
    // Verify the filter is applied
    await expect(this.page).toHaveURL(/category=/);
  }

  /**
   * Navigate to a specific page
   */
  async goToPage(pageNumber: number) {
    await this.page.locator(`[data-testid="page-${pageNumber}"]`).click();
    
    // Wait for the page to load
    await this.page.waitForResponse(response => 
      response.url().includes(`page=${pageNumber}`) && response.status() === 200
    );
    
    // Verify the URL contains the page parameter
    await expect(this.page).toHaveURL(new RegExp(`page=${pageNumber}`));
  }

  /**
   * Get todo count
   */
  async getTodoCount(): Promise<number> {
    return await this.page.locator('[data-testid="todo-item"]').count();
  }

  /**
   * Get category count
   */
  async getCategoryCount(): Promise<number> {
    return await this.page.locator('[data-testid="category-item"]').count();
  }

  /**
   * Check if element exists
   */
  async elementExists(selector: string): Promise<boolean> {
    return await this.page.locator(selector).count() > 0;
  }

  /**
   * Wait for element to be visible
   */
  async waitForElement(selector: string, timeout = 5000) {
    await this.page.waitForSelector(selector, { state: 'visible', timeout });
  }

  /**
   * Take a screenshot
   */
  async takeScreenshot(name: string) {
    await this.page.screenshot({ path: `test-results/screenshots/${name}.png` });
  }
} 