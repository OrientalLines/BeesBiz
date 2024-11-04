<script lang="ts">
  import { fly, fade } from 'svelte/transition';
  import { toasts } from '$lib/stores/toast';
  import { CheckCircle, XCircle, Info, AlertTriangle, X } from 'lucide-svelte';

  const icons = {
    success: CheckCircle,
    error: XCircle,
    info: Info,
    warning: AlertTriangle
  };

  const colors = {
    success: 'bg-green-100 text-green-800 border-green-200',
    error: 'bg-red-100 text-red-800 border-red-200',
    info: 'bg-blue-100 text-blue-800 border-blue-200',
    warning: 'bg-yellow-100 text-yellow-800 border-yellow-200'
  };
</script>

<div class="fixed top-4 right-4 z-50 space-y-2">
  {#each $toasts as toast (toast.id)}
    <div
      class="flex items-center gap-2 px-4 py-3 rounded-lg border shadow-sm {colors[toast.type]}"
      in:fly={{ x: 20, duration: 300 }}
      out:fade={{ duration: 200 }}
    >
      <svelte:component this={icons[toast.type]} class="w-5 h-5" />
      <p class="text-sm font-medium">{toast.message}</p>
      <button
        class="ml-2 hover:opacity-70"
        on:click={() => toasts.remove(toast.id)}
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  {/each}
</div> 