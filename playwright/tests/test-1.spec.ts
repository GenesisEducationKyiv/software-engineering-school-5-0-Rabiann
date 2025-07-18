import { test, expect } from '@playwright/test';

test('test', async ({ page }) => {
  const randomMail = crypto.randomUUID().split('-').join('') + "@mail.com"
  await page.goto('http://localhost:8000/');
  await expect(page).toHaveTitle(/Subscription Form/)
  await page.getByRole('textbox', { name: 'Enter email' }).click();
  await page.getByRole('textbox', { name: 'Enter email' }).fill(randomMail);
  await expect(page.getByRole('textbox', { name: 'Enter email' })).toHaveValue(randomMail)
  await page.getByRole('textbox', { name: 'Enter email' }).press('Tab');
  await page.getByRole('textbox', { name: 'Enter city' }).fill('kyiv');
  await expect(page.getByRole('textbox', { name: 'Enter city' })).toHaveValue('kyiv')
  await page.getByRole('radio', { name: 'Hourly' }).check();
  await page.getByRole('radio', { name: 'Daily' }).check();
  await expect(page.getByRole('radio', { name: 'Daily' })).toBeChecked()
  await page.getByRole('button', { name: 'Subscribe' }).click();
  await expect(page).toHaveTitle(/Confirmation Request/)
});