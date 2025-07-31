import { chromium, FullConfig } from '@playwright/test';

async function globalSetup(config: FullConfig) {
  const { baseURL } = config.projects[0].use;
  
  console.log('🌐 Starting global setup...');
  console.log(`📡 Base URL: ${baseURL}`);
  
  // Wait for the application to be ready
  const browser = await chromium.launch();
  const page = await browser.newPage();
  
  try {
    // Wait for the application to be available
    await page.goto(baseURL!);
    await page.waitForLoadState('networkidle');
    console.log('✅ Application is ready for testing');
  } catch (error) {
    console.error('❌ Failed to connect to application:', error);
    throw error;
  } finally {
    await browser.close();
  }
  
  console.log('🎯 Global setup completed');
}

export default globalSetup; 