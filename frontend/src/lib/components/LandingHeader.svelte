<script lang="ts">
  interface Props {
    title: string;
    tagline: string;
  }

  let { title, tagline }: Props = $props();
</script>

<div class="text-center space-y-2">
  <h1 class="text-5xl font-black tracking-tight text-surface-900-50">
    {title}
  </h1>
  <p class="text-xl font-medium text-surface-600-400 tracking-wide opacity-80">
    {tagline}
  </p>
</div>
