<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<script>
		// Force dark mode always
		document.documentElement.setAttribute("class", "scheme-dark");
	</script>
</svelte:head>

{@render children()}
