import type { Handle, HandleServerError } from '@sveltejs/kit';
import { AxiosError } from 'axios';
import jwt from 'jsonwebtoken';

export const handle: Handle = async ({ event, resolve }) => {
	const accessToken = event.cookies.get('access_token');

	let isSignedIn: boolean = false;
	let isAdmin: boolean = false;

	if (accessToken) {
		const jwtPayload = jwt.decode(accessToken, { json: true });
		if (jwtPayload?.exp && jwtPayload.exp * 1000 > Date.now()) {
			isSignedIn = true;
			isAdmin = jwtPayload?.isAdmin || false;
		}
	}

	if (event.url.pathname.startsWith('/settings') && !event.url.pathname.startsWith('/login')) {
		if (!isSignedIn) {
			return new Response(null, {
				status: 302,
				headers: { location: '/login' }
			});
		}
	}

	if (event.url.pathname.startsWith('/login') && isSignedIn) {
		return new Response(null, {
			status: 302,
			headers: { location: '/settings' }
		});
	}

	if (event.url.pathname.startsWith('/settings/admin') && !isAdmin) {
		return new Response(null, {
			status: 302,
			headers: { location: '/settings' }
		});
	}

	const response = await resolve(event);
	return response;
};

export const handleError: HandleServerError = async ({ error, message, status }) => {
	if (error instanceof AxiosError) {
		message = error.response?.data.error || message;
		status = error.response?.status || status;
		console.error(
			`Axios error: ${error.request.path} - ${error.response?.data.error ?? error.message}`
		);
	} else {
		console.error(error);
	}

	return {
		message,
		status
	};
};
