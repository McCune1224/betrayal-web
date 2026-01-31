<script lang="ts">
  interface Props {
    type?: 'error' | 'success' | 'info' | 'warning';
    message: string;
    className?: string;
  }

  let { type = 'error', message, className = '' }: Props = $props();

  const typeStyles = {
    error: 'bg-error-500/10 border-error-500/30 text-error-700-300',
    success: 'bg-success-500/10 border-success-500/30 text-success-700-300',
    info: 'bg-primary-500/10 border-primary-500/30 text-primary-700-300',
    warning: 'bg-warning-500/10 border-warning-500/30 text-warning-700-300'
  };
</script>

<div class="p-4 border-2 text-sm shadow-md {typeStyles[type]} {className}">
  {message}
</div>
