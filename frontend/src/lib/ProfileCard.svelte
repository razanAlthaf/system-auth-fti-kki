<script>
  import RoleBadge from './RoleBadge.svelte';

  export let profile;

  $: initials = profile.username.slice(0, 2).toUpperCase();
  $: joinedAt = new Date(profile.created_at).toLocaleDateString('id-ID', { dateStyle: 'long' });
  $: details = [
    { label: 'Email', value: profile.email },
    { label: 'Bergabung', value: joinedAt },
    { label: 'ID pengguna', value: profile.id, mono: true }
  ];
</script>

<section aria-labelledby="profile-title" class="rounded-lg border border-slate-200 bg-white p-6">
  <div class="flex items-center gap-4">
    <div
      class="flex h-14 w-14 shrink-0 items-center justify-center rounded-full bg-blue-800 text-lg font-semibold text-white"
      aria-hidden="true"
    >
      {initials}
    </div>
    <div class="min-w-0">
      <h2 id="profile-title" class="truncate text-xl font-semibold text-slate-900">{profile.username}</h2>
      {#if profile.role}<RoleBadge role={profile.role} />{/if}
    </div>
  </div>

  <dl class="mt-6 divide-y divide-slate-100 text-sm">
    {#each details as { label, value, mono }}
      <div class="flex flex-col gap-1 py-3 sm:flex-row sm:gap-6">
        <dt class="w-32 shrink-0 text-slate-500">{label}</dt>
        <dd class="min-w-0 break-all text-slate-900 {mono ? 'font-mono text-xs' : ''}">{value}</dd>
      </div>
    {/each}
  </dl>
</section>
