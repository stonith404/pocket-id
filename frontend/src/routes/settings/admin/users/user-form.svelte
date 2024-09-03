<script lang="ts">
	import FormInput from '$lib/components/form-input.svelte';
	import { Button } from '$lib/components/ui/button';
	import { Checkbox } from '$lib/components/ui/checkbox';
	import { Label } from '$lib/components/ui/label';
	import type { UserCreate } from '$lib/types/user.type';
	import { createForm } from '$lib/utils/form-util';
	import { z } from 'zod';

	let {
		callback,
		existingUser
	}: {
		existingUser?: UserCreate;
		callback: (user: UserCreate) => Promise<boolean>;
	} = $props();

	let isLoading = $state(false);

	const user = {
		firstName: existingUser?.firstName || '',
		lastName: existingUser?.lastName || '',
		email: existingUser?.email || '',
		username: existingUser?.username || '',
		isAdmin: existingUser?.isAdmin || false
	};

	const formSchema = z.object({
		firstName: z.string().min(2).max(30),
		lastName: z.string().min(2).max(30),
		username: z
			.string()
			.min(2)
			.max(30)
			.regex(
				/^[a-z0-9_@.-]+$/,
				"Username can only contain lowercase letters, numbers, underscores, dots, hyphens, and '@' symbols"
			),
		email: z.string().email(),
		isAdmin: z.boolean()
	});
	type FormSchema = typeof formSchema;

	const { inputs, ...form } = createForm<FormSchema>(formSchema, user);
	async function onSubmit() {
		const data = form.validate();
		if (!data) return;
		isLoading = true;
		const success = await callback(data);
		// Reset form if user was successfully created
		if (success && !existingUser) form.reset();
		isLoading = false;
	}
</script>

<form onsubmit={onSubmit}>
	<div class="flex flex-col gap-3 sm:flex-row">
		<div class="w-full">
			<FormInput label="Firstname" bind:input={$inputs.firstName} />
		</div>
		<div class="w-full">
			<FormInput label="Lastname" bind:input={$inputs.lastName} />
		</div>
	</div>
	<div class="mt-3 flex flex-col gap-3 sm:flex-row">
		<div class="w-full">
			<FormInput label="Email" bind:input={$inputs.email} />
		</div>
		<div class="w-full">
			<FormInput label="Username" bind:input={$inputs.username} />
		</div>
	</div>
	<div class="items-top mt-5 flex space-x-2">
		<Checkbox id="admin-privileges" bind:checked={$inputs.isAdmin.value} />
		<div class="grid gap-1.5 leading-none">
			<Label for="admin-privileges" class="mb-0 text-sm font-medium leading-none">
				Admin Privileges
			</Label>
			<p class="text-muted-foreground text-[0.8rem]">Admins have full access to the admin panel.</p>
		</div>
	</div>
	<div class="mt-5 flex justify-end">
		<Button {isLoading} type="submit">Save</Button>
	</div>
</form>
