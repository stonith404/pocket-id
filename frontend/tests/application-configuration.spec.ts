import test, { expect } from '@playwright/test';
import { cleanupBackend } from './utils/cleanup.util';

test.beforeEach(cleanupBackend);

test('Update general configuration', async ({ page }) => {
	await page.goto('/settings/admin/application-configuration');

	await page.getByLabel('Name').fill('Updated Name');
	await page.getByLabel('Session Duration').fill('30');
	await page.getByRole('button', { name: 'Save' }).first().click();

	await expect(page.getByRole('status')).toHaveText(
		'Application configuration updated successfully'
	);
	await expect(page.getByTestId('application-name')).toHaveText('Updated Name');

	await page.reload();

	await expect(page.getByLabel('Name')).toHaveValue('Updated Name');
	await expect(page.getByLabel('Session Duration')).toHaveValue('30');
});

test('Update application images', async ({ page }) => {
	await page.goto('/settings/admin/application-configuration');

	await page.getByLabel('Favicon').setInputFiles('tests/assets/w3-schools-favicon.ico');
	await page.getByLabel('Logo').setInputFiles('tests/assets/pingvin-share-logo.png');
	await page.getByLabel('Background Image').setInputFiles('tests/assets/clouds.jpg');
	await page.getByRole('button', { name: 'Save' }).nth(1).click();

	await expect(page.getByRole('status')).toHaveText('Images updated successfully');

	await page.request
		.get('/api/application-configuration/favicon')
		.then((res) => expect.soft(res.status()).toBe(200));
	await page.request
		.get('/api/application-configuration/logo')
		.then((res) => expect.soft(res.status()).toBe(200));

	await page.request
		.get('/api/application-configuration/background-image')
		.then((res) => expect.soft(res.status()).toBe(200));
});
