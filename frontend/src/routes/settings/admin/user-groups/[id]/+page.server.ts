import UserGroupService from '$lib/services/user-group-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const userGroupService = new UserGroupService(cookies.get('access_token'));
	const userGroup = await userGroupService.get(params.id);

	return { userGroup };
};
