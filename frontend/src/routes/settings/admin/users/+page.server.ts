import UserService from '$lib/services/user-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const userService = new UserService(cookies.get('access_token'));
	const users = await userService.list();
	return users;
};
