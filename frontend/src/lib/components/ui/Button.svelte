<script lang="ts">
  interface Props {
    variant?: 'primary' | 'secondary' | 'success' | 'error' | 'ghost';
    size?: 'sm' | 'md' | 'lg';
    disabled?: boolean;
    onclick?: () => void;
    type?: 'button' | 'submit' | 'reset';
    className?: string;
    children?: import('svelte').Snippet;
  }

  let {
    variant = 'primary',
    size = 'md',
    disabled = false,
    onclick,
    type = 'button',
    className = '',
    children
  }: Props = $props();

  const variantStyles = {
    primary: 'bg-primary-500 hover:bg-primary-400 border-primary-600 shadow-lg shadow-primary-500/30 text-white',
    secondary: 'bg-secondary-500 hover:bg-secondary-400 border-secondary-600 shadow-lg shadow-secondary-500/30 text-white',
    success: 'bg-success-500 hover:bg-success-400 border-success-600 shadow-lg shadow-success-500/30 text-white',
    error: 'bg-error-500 hover:bg-error-400 border-error-600 shadow-lg shadow-error-500/30 text-white',
    ghost: 'bg-white/20 hover:bg-white/30 border-transparent shadow text-white'
  };

  const sizeStyles = {
    sm: 'px-3 py-1 text-sm',
    md: 'px-4 py-2',
    lg: 'px-6 py-3 text-lg'
  };
</script>

<button
  {type}
  class="font-semibold transition-colors border-2 disabled:opacity-50 disabled:cursor-not-allowed {variantStyles[variant]} {sizeStyles[size]} {className}"
  {disabled}
  onclick={onclick}
>
  {@render children?.()}
</button>
