<script lang="ts">
  interface Props {
    name: string;
    size?: 'sm' | 'md' | 'lg';
    variant?: 'primary' | 'secondary' | 'surface';
  }

  let { name, size = 'md', variant = 'primary' }: Props = $props();

  function getInitial(name: string): string {
    return (name || '?').charAt(0).toUpperCase();
  }

  const sizeStyles = {
    sm: 'w-6 h-6 text-xs',
    md: 'w-8 h-8 text-sm',
    lg: 'w-12 h-12 text-lg'
  };

  const variantStyles = {
    primary: 'bg-primary-500 text-white',
    secondary: 'bg-secondary-500 text-white',
    surface: 'bg-surface-400-600 text-white'
  };
</script>

<div class="flex-shrink-0 flex items-center justify-center font-bold shadow {sizeStyles[size]} {variantStyles[variant]}">
  {getInitial(name)}
</div>
