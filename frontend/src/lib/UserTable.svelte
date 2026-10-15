<script>
  import { createEventDispatcher } from 'svelte';
  import RoleBadge from './RoleBadge.svelte';
  import Button from './Button.svelte';

  export let users = [];

  const dispatch = createEventDispatcher();
  let pendingUserId = null;

  function confirmForceLogout(userId) {
    pendingUserId = null;
    dispatch('forcelogout', userId);
  }
</script>

<section aria-labelledby="users-title" class="rounded-lg border border-slate-200 bg-white">
  <div class="flex items-baseline justify-between px-6 py-4">
    <h2 id="users-title" class="text-lg font-semibold text-slate-900">Pengguna terdaftar</h2>
    <span class="text-sm text-slate-500">{users.length} akun</span>
  </div>

  <div class="overflow-x-auto border-t border-slate-100">
    <table class="w-full min-w-[34rem] text-left text-sm">
      <thead class="bg-slate-50 text-slate-500">
        <tr>
          <th scope="col" class="px-6 py-2.5 font-medium">Username</th>
          <th scope="col" class="px-6 py-2.5 font-medium">Email</th>
          <th scope="col" class="px-6 py-2.5 font-medium">Peran</th>
          <th scope="col" class="px-6 py-2.5 text-right font-medium">Aksi</th>
        </tr>
      </thead>
      <tbody class="divide-y divide-slate-100">
        {#each users as user (user.id)}
          <tr class="hover:bg-slate-50/60">
            <td class="px-6 py-3 font-medium text-slate-900">{user.username}</td>
            <td class="px-6 py-3 text-slate-600">{user.email}</td>
            <td class="px-6 py-3"><RoleBadge role={user.role} /></td>
            <td class="px-6 py-3 text-right">
              {#if user.role === 'admin'}
                <span class="text-slate-300" aria-label="Tidak ada aksi">–</span>
              {:else if pendingUserId === user.id}
                <span class="inline-flex items-center gap-2">
                  <span class="text-xs text-slate-600">Keluarkan {user.username}?</span>
                  <Button size="sm" variant="danger" on:click={() => confirmForceLogout(user.id)}>Ya</Button>
                  <Button size="sm" variant="outline" on:click={() => (pendingUserId = null)}>Batal</Button>
                </span>
              {:else}
                <Button size="sm" variant="dangerOutline" on:click={() => (pendingUserId = user.id)}>
                  Paksa keluar
                </Button>
              {/if}
            </td>
          </tr>
        {:else}
          <tr>
            <td colspan="4" class="px-6 py-10 text-center text-slate-500">Belum ada pengguna terdaftar.</td>
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</section>
