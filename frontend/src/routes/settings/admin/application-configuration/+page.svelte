<script lang="ts">
	import { env } from '$env/dynamic/public';
	import CollapsibleCard from '$lib/components/collapsible-card.svelte';
	import AppConfigService from '$lib/services/app-config-service';
	import appConfigStore from '$lib/stores/application-configuration-store';
	import type { AllAppConfig } from '$lib/types/application-configuration';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { toast } from 'svelte-sonner';
	import AppConfigEmailForm from './forms/app-config-email-form.svelte';
	import AppConfigGeneralForm from './forms/app-config-general-form.svelte';
	import AppConfigLdapForm from './forms/app-config-ldap-form.svelte';
	import UpdateApplicationImages from './update-application-images.svelte';

	let { data } = $props();
	let appConfig = $state(data.appConfig);

	const uiConfigDisabled = env.PUBLIC_UI_CONFIG_DISABLED === 'true';
	const appConfigService = new AppConfigService();

	async function updateAppConfig(updatedAppConfig: Partial<AllAppConfig>) {
		appConfig = await appConfigService
			.update({
				...appConfig,
				...updatedAppConfig
			})
			.catch((e) => {
				axiosErrorToast(e);
				throw e;
			});
		await appConfigStore.reload();
	}

	async function updateImages(
		logoLight: File | null,
		logoDark: File | null,
		backgroundImage: File | null,
		favicon: File | null
	) {
		const faviconPromise = favicon ? appConfigService.updateFavicon(favicon) : Promise.resolve();
		const lightLogoPromise = logoLight
			? appConfigService.updateLogo(logoLight, true)
			: Promise.resolve();
		const darkLogoPromise = logoDark
			? appConfigService.updateLogo(logoDark, false)
			: Promise.resolve();
		const backgroundImagePromise = backgroundImage
			? appConfigService.updateBackgroundImage(backgroundImage)
			: Promise.resolve();

		await Promise.all([lightLogoPromise, darkLogoPromise, backgroundImagePromise, faviconPromise])
			.then(() => toast.success('Images updated successfully'))
			.catch(axiosErrorToast);
	}
</script>

<svelte:head>
	<title>Application Configuration</title>
</svelte:head>

<fieldset class="flex flex-col gap-5" disabled={uiConfigDisabled}>
	<CollapsibleCard id="application-configuration-general" title="General" defaultExpanded>
		<AppConfigGeneralForm {appConfig} callback={updateAppConfig} />
	</CollapsibleCard>

	<CollapsibleCard
		id="application-configuration-email"
		title="Email"
		description="Enable email notifications to alert users when a login is detected from a new device or
			location."
	>
		<AppConfigEmailForm {appConfig} callback={updateAppConfig} />
	</CollapsibleCard>

	<CollapsibleCard
		id="application-configuration-ldap"
		title="LDAP"
		description="Configure LDAP settings to sync users and groups from an LDAP server."
	>
		<AppConfigLdapForm {appConfig} callback={updateAppConfig} />
	</CollapsibleCard>
</fieldset>

<CollapsibleCard id="application-configuration-images" title="Images">
	<UpdateApplicationImages callback={updateImages} />
</CollapsibleCard>
