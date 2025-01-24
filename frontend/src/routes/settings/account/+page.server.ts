import { ACCESS_TOKEN_COOKIE_NAME } from '$lib/constants';
import UserService from '$lib/services/user-service';
import WebAuthnService from '$lib/services/webauthn-service';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ cookies }) => {
	const webauthnService = new WebAuthnService(cookies.get(ACCESS_TOKEN_COOKIE_NAME));
	const userService = new UserService(cookies.get(ACCESS_TOKEN_COOKIE_NAME));
	const account = await userService.getCurrent();
	const passkeys = await webauthnService.listCredentials();
	return {
		account,
		passkeys
	};
};
