<script>
  import { createEventDispatcher } from 'svelte';
  import { api, ApiError } from './api.js';
  import Field from './Field.svelte';
  import Button from './Button.svelte';

  export let notice = '';

  const dispatch = createEventDispatcher();
  const NETWORK_ERROR = 'Tidak dapat terhubung ke server. Coba lagi sebentar.';

  let mode = 'login';
  let submitting = false;
  let error = '';
  let info = notice;

  let loginForm = { identifier: '', password: '' };
  let registerForm = { email: '', username: '', password: '' };

  $: isLogin = mode === 'login';

  function switchMode(next) {
    mode = next;
    error = '';
    info = '';
  }

  async function run(task, failureMessage) {
    submitting = true;
    error = '';
    info = '';
    try {
      await task();
    } catch (err) {
      error = err instanceof ApiError ? failureMessage : NETWORK_ERROR;
    } finally {
      submitting = false;
    }
  }

  function handleLogin() {
    return run(async () => {
      const { token, role } = await api.login(loginForm);
      dispatch('authenticated', { token, role });
    }, 'Email/username atau kata sandi tidak cocok.');
  }

  function handleRegister() {
    return run(async () => {
      await api.register(registerForm);
      loginForm = { identifier: registerForm.email, password: '' };
      registerForm = { email: '', username: '', password: '' };
      mode = 'login';
      info = 'Akun berhasil dibuat. Silakan masuk.';
    }, 'Email atau username sudah dipakai akun lain.');
  }
</script>

<div class="grid min-h-screen lg:grid-cols-[5fr_6fr]">
  <aside class="flex flex-col justify-between bg-slate-900 px-8 py-10 text-slate-100 lg:px-14 lg:py-16">
    <p class="text-sm font-semibold tracking-wide text-slate-300">FTI</p>
    <div class="py-10 lg:py-0">
      <h1 class="max-w-sm text-3xl font-semibold leading-tight lg:text-4xl">
        Sistem Authentication FTI
      </h1>
      <p class="mt-4 max-w-sm leading-relaxed text-slate-400">
        Masuk untuk melihat profil Anda. Administrator dapat mengelola akun dan sesi pengguna dari sini.
      </p>
    </div>
    <p class="hidden text-xs text-slate-500 lg:block">Fakultas Teknologi Informasi</p>
  </aside>

  <main class="flex items-center justify-center bg-slate-50 px-6 py-12">
    <div class="w-full max-w-sm">
      <div role="tablist" class="mb-8 flex border-b border-slate-200">
        {#each [['login', 'Masuk'], ['register', 'Daftar']] as [key, text]}
          <button
            type="button"
            role="tab"
            aria-selected={mode === key}
            on:click={() => switchMode(key)}
            class="-mb-px border-b-2 px-4 pb-3 text-sm font-semibold transition-colors
                   focus:outline-none focus-visible:text-blue-800
                   {mode === key ? 'border-blue-800 text-blue-900' : 'border-transparent text-slate-500 hover:text-slate-800'}"
          >
            {text}
          </button>
        {/each}
      </div>

      <div aria-live="polite">
        {#if error}
          <p role="alert" class="mb-5 rounded-md border border-red-200 bg-red-50 px-3 py-2.5 text-sm text-red-800">
            {error}
          </p>
        {/if}
        {#if info}
          <p class="mb-5 rounded-md border border-emerald-200 bg-emerald-50 px-3 py-2.5 text-sm text-emerald-800">
            {info}
          </p>
        {/if}
      </div>

      {#if isLogin}
        <form on:submit|preventDefault={handleLogin} class="space-y-5">
          <Field id="identifier" label="Email atau username" autocomplete="username" bind:value={loginForm.identifier} />
          <Field id="password" label="Kata sandi" type="password" autocomplete="current-password" bind:value={loginForm.password} />
          <Button type="submit" full loading={submitting}>
            {submitting ? 'Memeriksa…' : 'Masuk'}
          </Button>
        </form>
      {:else}
        <form on:submit|preventDefault={handleRegister} class="space-y-5">
          <Field id="reg-email" label="Email" type="email" autocomplete="email" bind:value={registerForm.email} />
          <Field id="reg-username" label="Username" autocomplete="username" bind:value={registerForm.username} />
          <Field id="reg-password" label="Kata sandi" type="password" autocomplete="new-password" bind:value={registerForm.password} />
          <Button type="submit" full loading={submitting}>
            {submitting ? 'Membuat akun…' : 'Buat akun'}
          </Button>
        </form>
      {/if}
    </div>
  </main>
</div>
