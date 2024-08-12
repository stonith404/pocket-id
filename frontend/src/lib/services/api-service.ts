import { browser } from '$app/environment';
import { env } from '$env/dynamic/public';
import axios from 'axios';

abstract class APIService {
	baseURL: string = '/api';
	api = axios.create({
		withCredentials: true
	});

	constructor(accessToken?: string) {
		if (accessToken) {
			this.api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
		} else {
			this.api.defaults.baseURL = '/api';
		}
		if (!browser) {
			this.api.defaults.baseURL = (env.PUBLIC_APP_URL ?? 'http://localhost') + '/api';
		}
	}
}

export default APIService;
