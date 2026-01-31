<script lang="ts">
  interface Props {
    className?: string;
    children?: import('svelte').Snippet;
  }

  let { className = '', children }: Props = $props();
</script>

<div class="bg-surface-100-900 border-2 border-surface-300-700 shadow-2xl overflow-hidden backdrop-blur-sm {className}">
  {@render children?.()}
</div>
