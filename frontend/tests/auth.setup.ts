import { test as setup } from '@playwright/test';
import passkeyUtil from './utils/passkey.util';
import { cleanupBackend } from './utils/cleanup.util';

const authFile = 'tests/.auth/user.json';

setup('authenticate', async ({ page }) => {
	await cleanupBackend();
	await page.goto('/login');

	await (await passkeyUtil.init(page)).addPasskey();

	await page.getByRole('button', { name: 'Authenticate' }).click();
	await page.waitForURL('/settings/account');

	await page.context().storageState({ path: authFile });
});
