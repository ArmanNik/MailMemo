<div class="mm-container flex h-full min-h-[100vh] w-full flex-col lg:grid lg:grid-cols-2">
	<div class="mm-container-local hidden h-full w-full items-center justify-center lg:flex">
		<!-- <img
			src="/images/phone.png"
			alt="placeholder"
			width="1920"
			height="1080"
			class="h-full w-full object-cover"
		/> -->
		<img src="/logotype.svg" width="300" alt="" />
	</div>
	<div class="mt-8 block px-8 text-center lg:hidden"><img src="/logotype.svg" alt="" /></div>
	<div class="flex h-full grow flex-col items-center justify-between py-8">
		<div class="flex items-center justify-center py-12">
			<div class="mx-auto grid w-[350px] gap-6">
				<slot />
			</div>
		</div>
		<div class="pb-6 text-center text-sm sm:mt-auto">
			<p>
				<slot name="footer"></slot>
			</p>
		</div>
	</div>
</div>

<style>
	.mm-container-local {
		position: relative;
		width: 100%;
		height: 100%;
		overflow: hidden;
	}

	.mm-container-local::before {
		content: '';
		position: absolute;
		top: 50%;
		left: 50%;
		width: 100%;
		height: 100%;
		transform: translate(-50%, -50%);
		border-radius: 588px;
		opacity: 0.3;
		background: radial-gradient(49.55% 43.54% at 47% 50.69%, #cbc2ff 32%, #7c67fe 100%);
		filter: blur(150px);
		z-index: -1;
	}
</style>
