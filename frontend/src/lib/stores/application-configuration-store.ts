import ApplicationConfigurationService from '$lib/services/application-configuration-service';
import type { ApplicationConfiguration } from '$lib/types/application-configuration';
import { writable } from 'svelte/store';

const applicationConfigurationStore = writable<ApplicationConfiguration>();

const applicationConfigurationService = new ApplicationConfigurationService();

const reload = async () => {
	const applicationConfiguration = await applicationConfigurationService.list();
	applicationConfigurationStore.set(applicationConfiguration);
};

const set = (applicationConfiguration: ApplicationConfiguration) => {
	applicationConfigurationStore.set(applicationConfiguration);
}

export default {
	subscribe: applicationConfigurationStore.subscribe,
	reload,
	set
};
