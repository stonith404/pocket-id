import type { Paginated, PaginationRequest } from '$lib/types/pagination.type';
import type { User, UserCreate } from '$lib/types/user.type';
import APIService from './api-service';

export default class UserService extends APIService {
	async list(search?: string, pagination?: PaginationRequest) {
		const res = await this.api.get('/users', {
			params: {
				search,
				...pagination
			}
		});
		return res.data as Paginated<User>;
	}

	async get(id: string) {
		const res = await this.api.get(`/users/${id}`);
		return res.data as User;
	}

	async getCurrent() {
		const res = await this.api.get('/users/me');
		return res.data as User;
	}

	async create(user: UserCreate) {
		const res = await this.api.post('/users', user);
		return res.data as User;
	}

	async update(id: string, user: UserCreate) {
		const res = await this.api.put(`/users/${id}`, user);
		return res.data as User;
	}

	async updateCurrent(user: UserCreate) {
		const res = await this.api.put('/users/me', user);
		return res.data as User;
	}

	async remove(id: string) {
		await this.api.delete(`/users/${id}`);
	}

	async createOneTimeAccessToken(userId: string) {
		const res = await this.api.post(`/users/${userId}/one-time-access-token`, {
			userId,
			expiresAt: new Date(Date.now() + 1000 * 60 * 5).toISOString()
		});
		return res.data.token;
	}

	async exchangeOneTimeAccessToken(token: string) {
		const res = await this.api.post(`/one-time-access-token/${token}`);
		return res.data as User;
	}
}
