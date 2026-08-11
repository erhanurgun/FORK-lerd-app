<script lang="ts">
  import { openInEditor } from '$lib/editor';
  import type { QueryFrame } from '$lib/dumpsStream';
  import { m } from '../paraglide/messages.js';

  interface Props {
    src?: { file: string; line: number };
    trace?: QueryFrame[];
  }
  let { src, trace = [] }: Props = $props();
  let open = $state(false);
  let copied = $state(false);
  let copyTimer: ReturnType<typeof setTimeout> | null = null;
  async function copyPath() {
    if (!primary) return;
    try {
      await navigator.clipboard.writeText(`${primary.file}:${primary.line}`);
      copied = true;
      if (copyTimer) clearTimeout(copyTimer);
      copyTimer = setTimeout(() => (copied = false), 1500);
    } catch {
      /* clipboard unavailable; leave the path untouched */
    }
  }

  // The most useful single frame: first application frame, then innermost, then src.
  const primary = $derived(
    trace.find((f) => !f.file.includes('/vendor/')) ??
      trace[0] ??
      (src?.file ? { func: '', file: src.file, line: src.line } : undefined)
  );
</script>

{#if primary}
  <div class="text-gray-700 dark:text-gray-200">
    {#if primary.func}<span class="font-semibold">{primary.func}</span> · {/if}
    <span class="inline-flex items-center align-middle">
      <button
        type="button"
        class="font-mono text-lerd-red hover:underline break-all"
        onclick={() => openInEditor(primary.file, primary.line)}
        title={m.queries_openInEditor()}
      >{primary.file}:{primary.line}</button>
      <button
        type="button"
        class="shrink-0 px-2 flex items-center text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 hover:bg-gray-50 dark:hover:bg-white/5 border-l border-gray-100 dark:border-lerd-border/50 {copied ? 'text-emerald-600 dark:text-emerald-500' : ''}"
        onclick={copyPath}
        title={m.queries_copyPath()}
        aria-label={m.queries_copyPath()}
      >
      {#if copied}
        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24"><path d="M20 6L9 17l-5-5" /></svg>
      {:else}
        <svg class="w-3 h-3" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" viewBox="0 0 24 24"><rect x="9" y="9" width="13" height="13" rx="2" /><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" /></svg>
      {/if}
    </button>
    </span>
  </div>
{/if}
{#if trace.length > 1}
  <div>
    <button
      type="button"
      class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 underline"
      onclick={() => (open = !open)}
    >{open ? m.queries_hideTrace() : m.queries_details()}</button>
    {#if open}
      <ol class="font-mono space-y-0.5 mt-1">
        {#each trace as frame}
          {@const app = !frame.file.includes('/vendor/')}
          <li class={app ? 'text-gray-700 dark:text-gray-200' : 'text-gray-400 dark:text-gray-500'}>
            <span class={app ? 'font-semibold' : ''}>{frame.func}</span> ·
            <button
              type="button"
              class="hover:underline break-all {app ? 'text-lerd-red' : ''}"
              onclick={() => openInEditor(frame.file, frame.line)}
              title={m.queries_openInEditor()}
            >{frame.file}:{frame.line}</button>
          </li>
        {/each}
      </ol>
    {/if}
  </div>
{/if}
