import { browser } from '$app/environment';
import axios from 'axios';

abstract class APIService {
	api = axios.create({
		withCredentials: true
	});

	constructor(accessToken?: string) {
		if (accessToken) {
			this.api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`;
		}
		if (browser) {
			this.api.defaults.baseURL = '/api';
		} else {
			this.api.defaults.baseURL = process!.env!.INTERNAL_BACKEND_URL + '/api';
		}
	}
}

export default APIService;
