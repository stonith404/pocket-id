<script lang="ts">
	import * as Card from '$lib/components/ui/card';
	import ApplicationConfigurationService from '$lib/services/application-configuration-service';
	import applicationConfigurationStore from '$lib/stores/application-configuration-store';
	import type { AllApplicationConfiguration } from '$lib/types/application-configuration';
	import { axiosErrorToast } from '$lib/utils/error-util';
	import { toast } from 'svelte-sonner';
	import ApplicationConfigurationForm from './application-configuration-form.svelte';
	import UpdateApplicationImages from './update-application-images.svelte';

	let { data } = $props();
	let applicationConfiguration = $state(data.applicationConfiguration);

	const applicationConfigurationService = new ApplicationConfigurationService();

	async function updateConfiguration(configuration: AllApplicationConfiguration) {
		await applicationConfigurationService
			.update(configuration)
			.then(() => toast.success('Application configuration updated successfully'))
			.catch(axiosErrorToast);
		await applicationConfigurationStore.reload();
	}

	async function updateImages(
		logo: File | null,
		backgroundImage: File | null,
		favicon: File | null
	) {
		const faviconPromise = favicon
			? applicationConfigurationService.updateFavicon(favicon)
			: Promise.resolve();
		const logoPromise = logo ? applicationConfigurationService.updateLogo(logo) : Promise.resolve();
		const backgroundImagePromise = backgroundImage
			? applicationConfigurationService.updateBackgroundImage(backgroundImage)
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
		<ApplicationConfigurationForm {applicationConfiguration} callback={updateConfiguration} />
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
