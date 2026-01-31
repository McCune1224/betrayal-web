<script lang="ts">
  import Button from '$lib/components/ui/Button.svelte';

  interface Props {
    value: string;
    placeholder?: string;
    disabled?: boolean;
    onSend: () => void;
    onKeydown: (event: KeyboardEvent) => void;
  }

  let { value = $bindable(), placeholder = 'Type a message...', disabled = false, onSend, onKeydown }: Props = $props();
</script>

<div class="p-4 border-t border-surface-200-800 bg-surface-100-900">
  <div class="flex gap-2">
    <input
      type="text"
      class="flex-1 px-4 py-3 bg-surface-50-950 text-surface-900-50 border-2 border-surface-300-700 focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all placeholder:text-surface-400-600 shadow-inner"
      {placeholder}
      bind:value
      onkeydown={onKeydown}
      {disabled}
    />
    <Button onclick={onSend} disabled={!value.trim() || disabled}>
      Send
    </Button>
  </div>
</div>
