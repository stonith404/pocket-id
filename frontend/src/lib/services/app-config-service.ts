import type { AllAppConfig, AppConfigRawResponse } from '$lib/types/application-configuration';
import APIService from './api-service';

export default class AppConfigService extends APIService {
	async list(showAll = false) {
		let url = '/application-configuration';
		if (showAll) {
			url += '/all';
		}

		const { data } = await this.api.get<AppConfigRawResponse>(url);

		const appConfig: Partial<AllAppConfig> = {};
		data.forEach(({ key, value }) => {
			(appConfig as any)[key] = value;
		});

		return appConfig as AllAppConfig;
	}

	async update(appConfig: AllAppConfig) {
		const res = await this.api.put('/application-configuration', appConfig);
		return res.data as AllAppConfig;
	}

	async updateFavicon(favicon: File) {
		const formData = new FormData();
		formData.append('file', favicon!);

		await this.api.put(`/application-configuration/favicon`, formData);
	}

	async updateLogo(logo: File, light = true) {
		const formData = new FormData();
		formData.append('file', logo!);

		await this.api.put(`/application-configuration/logo`, formData, {
			params: { light }
		});
	}

	async updateBackgroundImage(backgroundImage: File) {
		const formData = new FormData();
		formData.append('file', backgroundImage!);

		await this.api.put(`/application-configuration/background-image`, formData);
	}
}
