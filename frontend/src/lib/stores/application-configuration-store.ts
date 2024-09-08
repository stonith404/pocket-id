import AppConfigService from '$lib/services/app-config-service';
import type { AppConfig } from '$lib/types/application-configuration';
import { writable } from 'svelte/store';

const appConfigStore = writable<AppConfig>();

const appConfigService = new AppConfigService();

const reload = async () => {
	const appConfig = await appConfigService.list();
	appConfigStore.set(appConfig);
};

const set = (appConfig: AppConfig) => {
	appConfigStore.set(appConfig);
};

export default {
	subscribe: appConfigStore.subscribe,
	reload,
	set
};
