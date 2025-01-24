import { ACCESS_TOKEN_COOKIE_NAME } from '$lib/constants';
import UserService from '$lib/services/user-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const userService = new UserService(cookies.get(ACCESS_TOKEN_COOKIE_NAME));
	const users = await userService.list();
	return users;
};
