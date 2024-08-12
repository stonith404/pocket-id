import UserService from '$lib/services/user-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, cookies }) => {
	const userService = new UserService(cookies.get('access_token'));
	const user = await userService.get(params.id);
	return user;
};
