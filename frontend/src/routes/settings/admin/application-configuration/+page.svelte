<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import AppConfigService from '$lib/services/app-config-service';
	import appConfigStore from '$lib/stores/application-configuration-store';
	import type { AllAppConfig } from '$lib/types/application-configuration';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { toast } from 'svelte-sonner';
	import AppConfigEmailForm from './forms/app-config-email-form.svelte';
	import AppConfigGeneralForm from './forms/app-config-general-form.svelte';
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
			.catch(axiosErrorToast);
		await appConfigStore.reload();
	}

	async function updateImages(
		logo: File | null,
		backgroundImage: File | null,
		favicon: File | null
	) {
		const faviconPromise = favicon ? appConfigService.updateFavicon(favicon) : Promise.resolve();
		const logoPromise = logo ? appConfigService.updateLogo(logo) : Promise.resolve();
		const backgroundImagePromise = backgroundImage
			? appConfigService.updateBackgroundImage(backgroundImage)
			: Promise.resolve();

		await Promise.all([logoPromise, backgroundImagePromise, faviconPromise])
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
			Email configuration is required for sending emails to users. If you disable email, users will
			not receive emails from the application.
		</Card.Description>
	</Card.Header>
	<Card.Content>
		<AppConfigEmailForm {appConfig} callback={updateAppConfig} />
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
