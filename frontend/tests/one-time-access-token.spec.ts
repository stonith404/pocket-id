import test, { expect } from '@playwright/test';
import { oneTimeAccessTokens } from './data';

// Disable authentication for these tests
test.use({ storageState: { cookies: [], origins: [] } });

test('Sign in with one time access token', async ({ page }) => {
	const token = oneTimeAccessTokens.filter((t) => !t.expired)[0];
	await page.goto(`/login/${token.token}`);

	await page.getByRole('button', { name: 'Continue' }).click();
	await page.waitForURL('/settings/account');
});

test('Sign in with expired one time access token fails', async ({ page }) => {
	const token = oneTimeAccessTokens.filter((t) => t.expired)[0];
	await page.goto(`/login/${token.token}`);

	await page.getByRole('button', { name: 'Continue' }).click();
	await expect(page.getByRole('paragraph')).toHaveText('Token is invalid or expired. Please try again.');
});
