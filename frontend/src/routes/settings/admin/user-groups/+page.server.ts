import UserGroupService from '$lib/services/user-group-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const userGroupService = new UserGroupService(cookies.get('access_token'));
	const userGroups = await userGroupService.list();
	return userGroups;
};
