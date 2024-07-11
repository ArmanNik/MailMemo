<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Separator } from '$lib/components/ui/separator';
	import Unauthenticated from '$lib/components/unauthenticated.svelte';
	import { account } from '$lib/sdk';
	import { OAuthProvider } from 'appwrite';

	let email = '';
	let password = '';

	async function handleSubmit() {
		try {
			await account.createEmailPasswordSession(email, password);
			await goto('/dashboard');
		} catch (error) {
			console.log(error);
		}
	}

	async function handleGitHubLogin() {
		try {
			account.createOAuth2Session(
				OAuthProvider.Github, // provider
				'https://mail-memo.vercel.app/dashboard', // redirect here on success
				'https://mail-memo.vercel.app/login', // redirect here on failure
				['repo', 'user']
			);
		} catch (error) {
			console.log(error);
		}
	}
</script>

<svelte:head>
	<title>Sign in - MailMemo</title>
</svelte:head>

<Unauthenticated>
	<div class="grid gap-2">
		<h1 class="font-header text-4xl">Sign in</h1>
		<p class="text-balance text-muted-foreground">Sign in to continue</p>
	</div>
	<form class="mt-9 grid gap-4" on:submit|preventDefault={handleSubmit}>
		<div class="grid gap-2">
			<Input
				class="frosted"
				id="email"
				type="email"
				placeholder="Email address"
				required
				bind:value={email}
			/>
		</div>
		<div class="grid gap-2">
			<Input
				class="frosted"
				id="password"
				type="password"
				required
				placeholder="Password"
				bind:value={password}
			/>
		</div>
		<Button type="submit" class="w-full">Sign in</Button>
		<div class=" grid grid-cols-3 items-center">
			<Separator />
			<p class="text-center text-sm text-muted-foreground">or</p>
			<Separator />
		</div>

		<Button variant="outline" class="frosted w-full" on:click={handleGitHubLogin}>
			<img src="/icons/github.svg" alt="github" />
			<span class="ml-2"> Login with GitHub </span>
		</Button>
	</form>

	<svelte:fragment slot="footer">
		Don't have an account?
		<a href="/register" class="underline"> Sign up </a>
	</svelte:fragment>
</Unauthenticated>
