import test, { expect } from '@playwright/test';
import { users } from './data';
import { cleanupBackend } from './utils/cleanup.util';

test.beforeEach(cleanupBackend);

test('Create user', async ({ page }) => {
	const user = users.steve;

	await page.goto('/settings/admin/users');

	await page.getByRole('button', { name: 'Add User' }).click();
	await page.getByLabel('Firstname').fill(user.firstname);
	await page.getByLabel('Lastname').fill(user.lastname);
	await page.getByLabel('Email').fill(user.email);
	await page.getByLabel('Username').fill(user.username);
	await page.getByRole('button', { name: 'Save' }).click();

	await expect(page.getByRole('row', { name: `${user.firstname} ${user.lastname}` })).toBeVisible();
	await expect(page.getByRole('status')).toHaveText('User created successfully');
});

test('Create user fails with already taken email', async ({ page }) => {
	const user = users.steve;

	await page.goto('/settings/admin/users');

	await page.getByRole('button', { name: 'Add User' }).click();
	await page.getByLabel('Firstname').fill(user.firstname);
	await page.getByLabel('Lastname').fill(user.lastname);
	await page.getByLabel('Email').fill(users.tim.email);
	await page.getByLabel('Username').fill(user.username);
	await page.getByRole('button', { name: 'Save' }).click();

	await expect(page.getByRole('status')).toHaveText('Email is already taken');
});

test('Create one time access token', async ({ page }) => {
	await page.goto('/settings/admin/users');

	await page
		.getByRole('row', { name: `${users.craig.firstname} ${users.craig.lastname}` })
		.getByRole('button')
		.click();
	await page.getByRole('menuitem', { name: 'One-time link' }).click();

	await expect(page.getByRole('textbox', { name: 'One Time Link' })).toHaveValue(
		/http:\/\/localhost\/login\/.*/
	);
});

test('Delete user', async ({ page }) => {
	await page.goto('/settings/admin/users');

	await page
		.getByRole('row', { name: `${users.craig.firstname} ${users.craig.lastname}` })
		.getByRole('button')
		.click();
	await page.getByRole('menuitem', { name: 'Delete' }).click();
	await page.getByRole('button', { name: 'Delete' }).click();

	await expect(page.getByRole('status')).toHaveText('User deleted successfully');
	await expect(
		page.getByRole('row', { name: `${users.craig.firstname} ${users.craig.lastname}` })
	).not.toBeVisible();
});

test('Update user', async ({ page }) => {
	const user = users.craig;

	await page.goto('/settings/admin/users');

	await page
		.getByRole('row', { name: `${user.firstname} ${user.lastname}` })
		.getByRole('button')
		.click();
	await page.getByRole('menuitem', { name: 'Edit' }).click();

	await page.getByLabel('Firstname').fill('Crack');
	await page.getByLabel('Lastname').fill('Apple');
	await page.getByLabel('Email').fill('crack.apple@test.com');
	await page.getByLabel('Username').fill('crack');
	await page.getByRole('button', { name: 'Save' }).click();

	await expect(page.getByRole('status')).toHaveText('User updated successfully');
});

test('Update user fails with already taken email', async ({ page }) => {
	const user = users.craig;

	await page.goto('/settings/admin/users');

	await page
		.getByRole('row', { name: `${user.firstname} ${user.lastname}` })
		.getByRole('button')
		.click();
	await page.getByRole('menuitem', { name: 'Edit' }).click();

	await page.getByLabel('Email').fill(users.tim.email);
	await page.getByRole('button', { name: 'Save' }).click();

	await expect(page.getByRole('status')).toHaveText('Email is already taken');
});

test('Update user fails with already taken username', async ({ page }) => {
	const user = users.craig;

	await page.goto('/settings/admin/users');

	await page
		.getByRole('row', { name: `${user.firstname} ${user.lastname}` })
		.getByRole('button')
		.click();
	await page.getByRole('menuitem', { name: 'Edit' }).click();

	await page.getByLabel('Username').fill(users.tim.username);
	await page.getByRole('button', { name: 'Save' }).click();

	await expect(page.getByRole('status')).toHaveText('Username is already taken');
});
