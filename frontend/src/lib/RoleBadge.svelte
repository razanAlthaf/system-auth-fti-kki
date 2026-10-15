<script>
  export let role;

  $: isAdmin = role === 'admin';
</script>

<span
  class="inline-block rounded-full px-2.5 py-0.5 text-xs font-medium capitalize
         {isAdmin ? 'bg-blue-100 text-blue-900' : 'bg-slate-100 text-slate-700'}"
>
  {role}
</span>
