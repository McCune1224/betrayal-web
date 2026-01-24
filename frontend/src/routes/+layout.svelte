<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import LightSwitch from '$lib/LightSwitch.svelte';

	let { children } = $props();
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>
<div style="position:fixed;top:1rem;right:1rem;z-index:1000;"><LightSwitch /></div>
{@render children()}
