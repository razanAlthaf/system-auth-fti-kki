<script>
  export let type = 'button';
  export let variant = 'primary';
  export let size = 'md';
  export let loading = false;
  export let disabled = false;
  export let full = false;

  const variants = {
    primary: 'bg-blue-800 text-white hover:bg-blue-900 focus-visible:ring-blue-800',
    outline: 'border border-slate-300 bg-white text-slate-700 hover:bg-slate-50 focus-visible:ring-slate-500',
    danger: 'bg-red-700 text-white hover:bg-red-800 focus-visible:ring-red-700',
    dangerOutline: 'border border-red-200 bg-white text-red-700 hover:bg-red-50 focus-visible:ring-red-600'
  };

  const sizes = {
    md: 'px-4 py-2.5 text-sm',
    sm: 'px-3 py-1.5 text-xs'
  };
</script>

<button
  {type}
  disabled={disabled || loading}
  on:click
  class="inline-flex items-center justify-center gap-2 rounded-md font-semibold transition-colors
         focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2
         disabled:cursor-not-allowed disabled:opacity-60
         {variants[variant]} {sizes[size]} {full ? 'w-full' : ''}"
>
  {#if loading}
    <span class="h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" aria-hidden="true"></span>
  {/if}
  <slot />
</button>
