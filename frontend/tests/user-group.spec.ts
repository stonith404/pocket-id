import test, { expect } from '@playwright/test';
import { userGroups, users } from './data';
import { cleanupBackend } from './utils/cleanup.util';

test.beforeEach(cleanupBackend);

test('Create user group', async ({ page, baseURL }) => {
	await page.goto('/settings/admin/user-groups');
	const group = userGroups.humanResources;

	await page.getByRole('button', { name: 'Add Group' }).click();
	await page.getByLabel('Friendly Name').fill(group.friendlyName);

	await page.getByRole('button', { name: 'Save' }).click();

	await expect(page.getByRole('status')).toHaveText('User group created successfully');

	const expectedRoute = new RegExp(`${baseURL}/settings/admin/user-groups/[a-f0-9-]+`);
	expect(page.url()).toMatch(expectedRoute);

	await expect(page.getByLabel('Friendly Name')).toHaveValue(group.friendlyName);
	await expect(page.getByLabel('Name', { exact: true })).toHaveValue(group.name);
});

test('Edit user group', async ({ page }) => {
	await page.goto('/settings/admin/user-groups');
	const group = userGroups.developers;

	await page.getByRole('row', { name: group.name }).getByRole('button').click();
	await page.getByRole('menuitem', { name: 'Edit' }).click();

	await page.getByLabel('Friendly Name').fill('Developers updated');

	await expect(page.getByLabel('Name', { exact: true })).toHaveValue(group.name);

	await page.getByLabel('Name', { exact: true }).fill('developers_updated');

	await page.getByRole('button', { name: 'Save' }).nth(0).click();

	await expect(page.getByRole('status')).toHaveText('User group updated successfully');
	await expect(page.getByLabel('Friendly Name')).toHaveValue('Developers updated');
	await expect(page.getByLabel('Name', { exact: true })).toHaveValue('developers_updated');
});

test('Update user group users', async ({ page }) => {
	const group = userGroups.designers;
	await page.goto(`/settings/admin/user-groups/${group.id}`);

	await page.getByRole('row', { name: users.tim.email }).getByRole('checkbox').click();
	await page.getByRole('row', { name: users.craig.email }).getByRole('checkbox').click();

	await page.getByRole('button', { name: 'Save' }).nth(1).click();

	await expect(page.getByRole('status')).toHaveText('Users updated successfully');

	await page.reload();

	await expect(
		page.getByRole('row', { name: users.tim.email }).getByRole('checkbox')
	).toHaveAttribute('data-state', 'unchecked');
	await expect(
		page.getByRole('row', { name: users.craig.email }).getByRole('checkbox')
	).toHaveAttribute('data-state', 'checked');
});

test('Delete user group', async ({ page }) => {
	const group = userGroups.developers;
	await page.goto('/settings/admin/user-groups');

	await page.getByRole('row', { name: group.name }).getByRole('button').click();
	await page.getByRole('menuitem', { name: 'Delete' }).click();
	await page.getByRole('button', { name: 'Delete' }).click();

	await expect(page.getByRole('status')).toHaveText('User group deleted successfully');
	await expect(page.getByRole('row', { name: group.name })).not.toBeVisible();
});
