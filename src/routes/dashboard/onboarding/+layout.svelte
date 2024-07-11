<script>
	import { step } from './store';
</script>

<div class="mt-4 flex h-full grow flex-col">
	<p class="text-muted-foreground">Step {$step}</p>
	<div class="flex h-full grow flex-col justify-between">
		<slot />
	</div>
</div>
