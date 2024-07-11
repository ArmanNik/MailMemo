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
