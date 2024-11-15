import type { AuthorizeResponse, OidcClient, OidcClientCreate } from '$lib/types/oidc.type';
import type { Paginated, PaginationRequest } from '$lib/types/pagination.type';
import APIService from './api-service';

class OidcService extends APIService {
	async authorize(clientId: string, scope: string, callbackURL: string, nonce?: string, codeChallenge?: string, codeChallengeMethod?: string) {
		const res = await this.api.post('/oidc/authorize', {
			scope,
			nonce,
			callbackURL,
			clientId,
			codeChallenge,
			codeChallengeMethod
		});

		return res.data as AuthorizeResponse;
	}

	async authorizeNewClient(clientId: string, scope: string, callbackURL: string, nonce?: string, codeChallenge?: string, codeChallengeMethod?: string) {
		const res = await this.api.post('/oidc/authorize/new-client', {
			scope,
			nonce,
			callbackURL,
			clientId,
			codeChallenge,
			codeChallengeMethod
		});

		return res.data as AuthorizeResponse;
	}

	async listClients(search?: string, pagination?: PaginationRequest) {
		const res = await this.api.get('/oidc/clients', {
			params: {
				search,
				...pagination
			}
		});
		return res.data as Paginated<OidcClient>;
	}

	async createClient(client: OidcClientCreate) {
		return (await this.api.post('/oidc/clients', client)).data as OidcClient;
	}

	async removeClient(id: string) {
		await this.api.delete(`/oidc/clients/${id}`);
	}

	async getClient(id: string) {
		return (await this.api.get(`/oidc/clients/${id}`)).data as OidcClient;
	}

	async updateClient(id: string, client: OidcClientCreate) {
		return (await this.api.put(`/oidc/clients/${id}`, client)).data as OidcClient;
	}

	async updateClientLogo(client: OidcClient, image: File | null) {
		if (client.hasLogo && !image) {
			await this.removeClientLogo(client.id);
			return;
		}
		if (!client.hasLogo && !image) {
			return;
		}

		const formData = new FormData();
		formData.append('file', image!);

		await this.api.post(`/oidc/clients/${client.id}/logo`, formData);
	}

	async removeClientLogo(id: string) {
		await this.api.delete(`/oidc/clients/${id}/logo`);
	}

	async createClientSecret(id: string) {
		return (await this.api.post(`/oidc/clients/${id}/secret`)).data.secret as string;
	}
}

export default OidcService;
