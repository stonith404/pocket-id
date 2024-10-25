import type { CustomClaim } from '$lib/types/custom-claim.type';
import type { User } from 'lucide-svelte';
import APIService from './api-service';

export default class CustomClaimService extends APIService {
	async getSuggestions() {
		const res = await this.api.get('/custom-claims/suggestions');
		return res.data as string[];
	}

	async updateUserCustomClaims(userId: string, claims: CustomClaim[]) {
		const res = await this.api.put(`/custom-claims/user/${userId}`, claims);
		return res.data as User;
	}
}
