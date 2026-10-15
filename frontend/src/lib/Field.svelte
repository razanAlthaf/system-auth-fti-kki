<script>
  export let id;
  export let label;
  export let value = '';
  export let type = 'text';
  export let autocomplete = undefined;
  export let hint = '';

  let passwordVisible = false;

  $: isPassword = type === 'password';
  $: inputType = isPassword && passwordVisible ? 'text' : type;

  function handleInput(event) {
    value = event.currentTarget.value;
  }
</script>

<div>
  <label for={id} class="mb-1.5 block text-sm font-medium text-slate-700">{label}</label>
  <div class="relative">
    <input
      {id}
      {value}
      {autocomplete}
      type={inputType}
      required
      on:input={handleInput}
      class="w-full rounded-md border border-slate-300 bg-white px-3 py-2.5 text-slate-900
             placeholder:text-slate-400 focus:border-blue-800 focus:outline-none focus:ring-2 focus:ring-blue-800/25
             {isPassword ? 'pr-24' : ''}"
    />
    {#if isPassword}
      <button
        type="button"
        on:click={() => (passwordVisible = !passwordVisible)}
        class="absolute inset-y-0 right-0 rounded-r-md px-3 text-xs font-medium text-slate-500
               hover:text-slate-800 focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-800"
      >
        {passwordVisible ? 'Sembunyikan' : 'Tampilkan'}
      </button>
    {/if}
  </div>
  {#if hint}
    <p class="mt-1 text-xs text-slate-500">{hint}</p>
  {/if}
</div>
