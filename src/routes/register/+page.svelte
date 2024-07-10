<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { account } from '$lib/sdk';
	import { ID, OAuthProvider } from 'appwrite';
	import { Separator } from '$lib/components/ui/separator';
	import { base } from '$app/paths';

	let email = '';
	let password = '';
	let name = '';

	async function handleSubmit() {
		try {
			await account.create(ID.unique(), email, password, name);
			await goto(`${base}/dashboard`);
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

<div class="flex h-full min-h-[100vh] w-full flex-col lg:grid lg:grid-cols-2">
	<div class="hidden bg-muted lg:block">
		<img
			src="/images/placeholder.svg"
			alt="placeholder"
			width="1920"
			height="1080"
			class="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
		/>
	</div>
	<div class="mt-8 block text-center lg:hidden">logo</div>
	<div class="flex h-full grow flex-col items-center justify-between py-8">
		<div class="flex items-center justify-center py-12">
			<div class="mx-auto grid w-[350px] gap-6">
				<div class="grid gap-2">
					<h1 class="text-3xl font-bold">Sign up</h1>
					<p class="text-balance text-muted-foreground">Sign up to continue</p>
				</div>
				<form class="grid gap-4" on:submit|preventDefault={handleSubmit}>
					<div class="grid gap-2">
						<Input id="name" type="text" placeholder="Full name" required bind:value={name} />
					</div>
					<div class="grid gap-2">
						<Input
							id="email"
							type="email"
							placeholder="Email address"
							required
							bind:value={email}
						/>
					</div>
					<div class="grid gap-2">
						<Input
							id="password"
							type="password"
							required
							placeholder="Password"
							bind:value={password}
						/>
					</div>
					<Button type="submit" class="w-full">Sign up</Button>
					<div class=" grid grid-cols-3 items-center">
						<Separator />
						<p class="text-center text-sm text-muted-foreground">or</p>
						<Separator />
					</div>

					<Button variant="outline" class="w-full" on:click={handleGitHubLogin}>
						Login with GitHub
					</Button>
				</form>
			</div>
		</div>
		<div class="text-center text-sm sm:mt-auto">
			Don't have an account?
			<a href="/register" class="underline"> Sign up </a>
		</div>
	</div>
</div>
