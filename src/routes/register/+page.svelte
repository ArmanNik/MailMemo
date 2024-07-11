<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { account } from '$lib/sdk';
	import { ID, OAuthProvider } from 'appwrite';
	import { Separator } from '$lib/components/ui/separator';
	import { base } from '$app/paths';
	import Unauthenticated from '$lib/components/unauthenticated.svelte';
	import { toast } from 'svelte-sonner';

	let email = '';
	let password = '';
	let name = '';
	let disabled = false;

	async function handleSubmit() {
		disabled = true;
		try {
			await account.create(ID.unique(), email, password, name);
			await account.createEmailPasswordSession(email, password);

			await goto(`${base}/dashboard`);
		} catch (error) {
			const message = (error as Error)?.message ?? error;
			toast(message);
			console.log(error);
		} finally {
			disabled = false;
		}
	}

	async function handleGitHubLogin() {
		try {
			account.createOAuth2Session(
				OAuthProvider.Github, // provider
				'https://app.mailmemo.site/dashboard', // redirect here on success
				'https://app.mailmemo.site/login', // redirect here on failure
				['repo', 'user']
			);
		} catch (error) {
			const message = (error as Error)?.message ?? error;
			toast(message);
			console.log(error);
		}
	}
</script>

<svelte:head>
	<title>Sign up - MailMemo</title>
</svelte:head>

<Unauthenticated>
	<div class="grid gap-2">
		<h1 class="font-header text-4xl">Sign up</h1>
		<p class="text-balance text-muted-foreground">Sign up to continue</p>
	</div>
	<form class="mt-9 grid gap-4" on:submit|preventDefault={handleSubmit}>
		<div class="grid gap-2">
			<Input
				id="name"
				type="text"
				placeholder="Full name"
				required
				class="frosted"
				bind:value={name}
			/>
		</div>
		<div class="grid gap-2">
			<Input
				id="email"
				type="email"
				placeholder="Email address"
				required
				class="frosted"
				bind:value={email}
			/>
		</div>
		<div class="grid gap-2">
			<Input
				id="password"
				type="password"
				required
				placeholder="Password"
				class="frosted"
				bind:value={password}
			/>
		</div>
		<Button type="submit" class="w-full" {disabled}>Sign up</Button>
		<div class=" grid grid-cols-3 items-center">
			<Separator />
			<p class="text-center text-sm text-muted-foreground">or</p>
			<Separator />
		</div>

		<Button variant="outline" class="frosted w-full" on:click={handleGitHubLogin}>
			<img src="/icons/github.svg" alt="github" />
			<span class="ml-2">Sign up with GitHub</span>
		</Button>
	</form>
	<svelte:fragment slot="footer">
		Already have an account?
		<a href="/login" class="underline">Sign in</a>
	</svelte:fragment>
</Unauthenticated>
