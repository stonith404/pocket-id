<script lang="ts">
	import * as Card from '$lib/components/ui/card';
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

	const appConfigService = new AppConfigService();

	async function updateAppConfig(updatedAppConfig: Partial<AllAppConfig>) {
		await appConfigService
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

<Card.Root>
	<Card.Header>
		<Card.Title>General</Card.Title>
	</Card.Header>
	<Card.Content>
		<AppConfigGeneralForm {appConfig} callback={updateAppConfig} />
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Email</Card.Title>
		<Card.Description>
			Enable email notifications to alert users when a login is detected from a new device or
			location.
		</Card.Description>
	</Card.Header>
	<Card.Content>
		<AppConfigEmailForm {appConfig} callback={updateAppConfig} />
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>LDAP</Card.Title>
		<Card.Description>
			Configure LDAP settings to sync users and groups from an LDAP server.
		</Card.Description>
	</Card.Header>
	<Card.Content>
		<AppConfigLdapForm {appConfig} callback={updateAppConfig} />
	</Card.Content>
</Card.Root>

<Card.Root>
	<Card.Header>
		<Card.Title>Images</Card.Title>
	</Card.Header>
	<Card.Content>
		<UpdateApplicationImages callback={updateImages} />
	</Card.Content>
</Card.Root>
