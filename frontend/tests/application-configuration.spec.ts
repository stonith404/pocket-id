import test, { expect } from '@playwright/test';
import { cleanupBackend } from './utils/cleanup.util';

test.beforeEach(cleanupBackend);

test('Update general configuration', async ({ page }) => {
	await page.goto('/settings/admin/application-configuration');

	await page.getByLabel('Application Name', { exact: true }).fill('Updated Name');
	await page.getByLabel('Session Duration').fill('30');
	await page.getByRole('button', { name: 'Save' }).first().click();

	await expect(page.getByRole('status')).toHaveText(
		'Application configuration updated successfully'
	);
	await expect(page.getByTestId('application-name')).toHaveText('Updated Name');

	await page.reload();

	await expect(page.getByLabel('Application Name', { exact: true })).toHaveValue('Updated Name');
	await expect(page.getByLabel('Session Duration')).toHaveValue('30');
});

test('Update email configuration', async ({ page }) => {
	await page.goto('/settings/admin/application-configuration');

	await page.getByLabel('SMTP Host').fill('smtp.gmail.com');
	await page.getByLabel('SMTP Port').fill('587');
	await page.getByLabel('SMTP User').fill('test@gmail.com');
	await page.getByLabel('SMTP Password').fill('password');
	await page.getByLabel('SMTP From').fill('test@gmail.com');
	await page.getByRole('button', { name: 'Enable' }).nth(0).click();
	await page.getByRole('status').click();

	await expect(page.getByRole('status')).toHaveText('Email configuration updated successfully');
	await expect(page.getByRole('button', { name: 'Disable' })).toBeVisible();

	await page.reload();

	await expect(page.getByLabel('SMTP Host')).toHaveValue('smtp.gmail.com');
	await expect(page.getByLabel('SMTP Port')).toHaveValue('587');
	await expect(page.getByLabel('SMTP User')).toHaveValue('test@gmail.com');
	await expect(page.getByLabel('SMTP Password')).toHaveValue('password');
	await expect(page.getByLabel('SMTP From')).toHaveValue('test@gmail.com');

	await page.getByRole('button', { name: 'Disable' }).click();

	await expect(page.getByRole('status')).toHaveText('Email disabled successfully');
});

test('Update application images', async ({ page }) => {
	await page.goto('/settings/admin/application-configuration');

	await page.getByLabel('Favicon').setInputFiles('tests/assets/w3-schools-favicon.ico');
	await page.getByLabel('Light Mode Logo').setInputFiles('tests/assets/pingvin-share-logo.png');
	await page.getByLabel('Dark Mode Logo').setInputFiles('tests/assets/nextcloud-logo.png');
	await page.getByLabel('Background Image').setInputFiles('tests/assets/clouds.jpg');
	await page.getByRole('button', { name: 'Save' }).nth(1).click();

	await expect(page.getByRole('status')).toHaveText('Images updated successfully');

	await page.request
		.get('/api/application-configuration/favicon')
		.then((res) => expect.soft(res.status()).toBe(200));
	await page.request
		.get('/api/application-configuration/logo?light=true')
		.then((res) => expect.soft(res.status()).toBe(200));
	await page.request
		.get('/api/application-configuration/logo?light=false')
		.then((res) => expect.soft(res.status()).toBe(200));
	await page.request
		.get('/api/application-configuration/background-image')
		.then((res) => expect.soft(res.status()).toBe(200));
});
