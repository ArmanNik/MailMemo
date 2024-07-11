<div class="mm-container flex h-full min-h-[100vh] w-full flex-col lg:grid lg:grid-cols-2">
	<div class="hidden lg:block">
		<img
			src="/images/phone.png"
			alt="placeholder"
			width="1920"
			height="1080"
			class="h-full w-full object-cover"
		/>
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
