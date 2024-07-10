<div class="mm-container flex h-full min-h-[100vh] w-full flex-col lg:grid lg:grid-cols-2">
	<div class="hidden bg-muted lg:block">
		<img
			src="/images/placeholder.svg"
			alt="placeholder"
			width="1920"
			height="1080"
			class="h-full w-full object-cover dark:brightness-[0.2] dark:grayscale"
		/>
	</div>
	<div class="mm-gradient"></div>
	<div class="mt-8 block text-center lg:hidden">logo</div>
	<div class="flex h-full grow flex-col items-center justify-between py-8">
		<div class="flex items-center justify-center py-12">
			<div class="mx-auto grid w-[350px] gap-6">
				<slot />
			</div>
		</div>
		<div class="text-center text-sm sm:mt-auto">
			<p>
				<slot name="footer"></slot>
			</p>
		</div>
	</div>
</div>

<style>
	.mm-container {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;
	}
	.mm-container::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 80%;
		width: 100%;
		height: 80%;
		transform: translateY(-50%);
		border-radius: 588px;
		opacity: 0.3;
		background: radial-gradient(49.55% 43.54% at 47% 50.69%, #cbc2ff 32%, #7c67fe 100%);
		filter: blur(150px);
	}
	.mm-gradient {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		border-radius: 588px;
		opacity: 0.3;
		background: radial-gradient(49.55% 43.54% at 47% 50.69%, #cbc2ff 32%, #7c67fe 100%);
		filter: blur(150px);
	}
</style>
