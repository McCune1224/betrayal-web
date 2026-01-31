<script lang="ts">
  interface Props {
    label?: string;
    value: string;
    placeholder?: string;
    type?: 'text' | 'password' | 'email';
    error?: string;
    maxlength?: number;
    disabled?: boolean;
    className?: string;
    oninput?: (value: string) => void;
    onkeydown?: (event: KeyboardEvent) => void;
  }

  let {
    label,
    value = $bindable(),
    placeholder = '',
    type = 'text',
    error,
    maxlength,
    disabled = false,
    className = '',
    oninput,
    onkeydown
  }: Props = $props();

  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    value = target.value;
    oninput?.(value);
  }
</script>

<label class="block space-y-2 {className}">
  {#if label}
    <span class="text-sm font-bold uppercase tracking-wider text-surface-600-400">
      {label}
    </span>
  {/if}
  <input
    {type}
    {placeholder}
    {maxlength}
    {disabled}
    {value}
    oninput={handleInput}
    onkeydown={onkeydown}
    class="w-full px-4 py-3 bg-surface-50-950 text-surface-900-50 border-2 {error ? 'border-error-500' : 'border-surface-300-700'} focus:border-primary-500 focus:ring-1 focus:ring-primary-500 outline-none transition-all placeholder:text-surface-400-600 shadow-inner"
  />
  {#if error}
    <span class="text-sm text-error-500">{error}</span>
  {/if}
</label>
