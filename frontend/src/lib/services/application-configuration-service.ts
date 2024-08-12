import type {
	AllApplicationConfiguration,
	ApplicationConfigurationRawResponse
} from '$lib/types/application-configuration';
import APIService from './api-service';

export default class ApplicationConfigurationService extends APIService {
	async list(showAll = false) {
		const { data } = await this.api.get<ApplicationConfigurationRawResponse>(
			'/application-configuration',
			{
				params: { showAll }
			}
		);

		const applicationConfiguration: Partial<AllApplicationConfiguration> = {};
		data.forEach(({ key, value }) => {
			(applicationConfiguration as any)[key] = value;
		});

		return applicationConfiguration as AllApplicationConfiguration;
	}

	async update(applicationConfiguration: AllApplicationConfiguration) {
		const res = await this.api.put('/application-configuration', applicationConfiguration);
		return res.data as AllApplicationConfiguration;
	}

	async updateFavicon(favicon: File) {
		const formData = new FormData();
		formData.append('file', favicon!);

		await this.api.put(`/application-configuration/favicon`, formData);
	}

	async updateLogo(logo: File) {
		const formData = new FormData();
		formData.append('file', logo!);

		await this.api.put(`/application-configuration/logo`, formData);
	}

	async updateBackgroundImage(backgroundImage: File) {
		const formData = new FormData();
		formData.append('file', backgroundImage!);

		await this.api.put(`/application-configuration/background-image`, formData);
	}
}
