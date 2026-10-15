<script>
  import { onMount } from 'svelte';
  import { api } from './lib/api.js';
  import { loadSession, saveSession, clearSession } from './lib/session.js';
  import AuthCard from './lib/AuthCard.svelte';
  import ProfileCard from './lib/ProfileCard.svelte';
  import UserTable from './lib/UserTable.svelte';
  import RoleBadge from './lib/RoleBadge.svelte';
  import Button from './lib/Button.svelte';

  const TOAST_DURATION_MS = 4000;
  const EMPTY_SESSION = { token: '', role: '' };

  let session = loadSession();
  let profile = null;
  let users = [];
  let toast = null;
  let toastTimer;
  let sessionNotice = '';
  let loggingOut = false;

  $: isLoggedIn = Boolean(session.token);
  $: isAdmin = session.role === 'admin';

  onMount(() => {
    if (isLoggedIn) loadDashboard();
  });

  function showToast(text, type = 'success') {
    clearTimeout(toastTimer);
    toast = { text, type };
    toastTimer = setTimeout(() => (toast = null), TOAST_DURATION_MS);
  }

  function endSession(notice = '') {
    clearSession();
    session = EMPTY_SESSION;
    profile = null;
    users = [];
    sessionNotice = notice;
  }

  async function loadProfile() {
    try {
      profile = await api.getProfile(session.token);
    } catch {
      endSession('Sesi Anda sudah berakhir. Silakan masuk kembali.');
    }
  }

  async function loadUsers() {
    try {
      users = await api.getUsers(session.token);
    } catch {
      showToast('Daftar pengguna gagal dimuat.', 'error');
    }
  }

  async function loadDashboard() {
    await loadProfile();
    if (isLoggedIn && isAdmin) await loadUsers();
  }

  function handleAuthenticated({ detail }) {
    session = detail;
    saveSession(detail);
    sessionNotice = '';
    loadDashboard();
  }

  async function handleLogout() {
    loggingOut = true;
    try {
      await api.logout(session.token);
    } catch {
      // Sesi lokal tetap dihapus meski server tidak terjangkau.
    } finally {
      loggingOut = false;
      endSession();
    }
  }

  async function handleForceLogout({ detail: userId }) {
    try {
      await api.forceLogout(session.token, userId);
      showToast('Sesi pengguna berhasil diakhiri.');
      await loadUsers();
    } catch {
      showToast('Sesi pengguna gagal diakhiri.', 'error');
    }
  }
</script>

{#if !isLoggedIn}
  <AuthCard notice={sessionNotice} on:authenticated={handleAuthenticated} />
{:else}
  <div class="min-h-screen bg-slate-50">
    <header class="border-b border-slate-200 bg-white">
      <div class="mx-auto flex max-w-4xl items-center justify-between gap-4 px-6 py-3">
        <span class="font-semibold text-slate-900">Authentication FTI</span>
        <div class="flex items-center gap-3">
          <RoleBadge role={session.role} />
          <Button variant="outline" size="sm" loading={loggingOut} on:click={handleLogout}>Keluar</Button>
        </div>
      </div>
    </header>

    <main class="mx-auto max-w-4xl space-y-6 px-6 py-8">
      {#if profile}
        <ProfileCard {profile} />
      {:else}
        <div class="h-52 animate-pulse rounded-lg border border-slate-200 bg-white" aria-busy="true"></div>
      {/if}

      {#if isAdmin}
        <UserTable {users} on:forcelogout={handleForceLogout} />
      {/if}
    </main>
  </div>
{/if}

{#if toast}
  <div
    role="status"
    class="fixed bottom-6 left-1/2 -translate-x-1/2 rounded-md px-4 py-2.5 text-sm font-medium text-white shadow-lg
           {toast.type === 'error' ? 'bg-red-700' : 'bg-slate-900'}"
  >
    {toast.text}
  </div>
{/if}

<style>
  :global(body) {
    font-family: 'Public Sans', system-ui, sans-serif;
    -webkit-font-smoothing: antialiased;
  }
</style>
