import { FullConfig } from '@playwright/test';

async function globalTeardown(config: FullConfig) {
  console.log('🧹 Starting global teardown...');
  
  // Clean up any test data or resources
  // This is where you would clean up the database, remove test files, etc.
  
  console.log('✅ Global teardown completed');
}

export default globalTeardown; 