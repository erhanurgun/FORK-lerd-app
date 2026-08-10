<script lang="ts">
  import { dashboardIconSvg } from '$lib/dashboardIcons';
  import { serviceMeta } from '$lib/serviceMeta';
  import { serviceIcons } from '$stores/serviceIcons';
  import { brandTintStyle } from '$lib/brandTint';
  import { tintFor, asCategory } from '$lib/presetCategories';

  interface Props {
    name: string;
    category?: string;
    icon?: string;
    color?: string;
    preset?: string;
    compact?: boolean;
  }
  let { name, category, icon, color, preset, compact = false }: Props = $props();

  // A caller holding the preset or service passes its declared values; one
  // holding only a name resolves them through the registry.
  const meta = $derived($serviceMeta.get(name));
  const resolvedCategory = $derived(asCategory(category ?? meta?.category));
  const resolvedIcon = $derived(icon ?? meta?.icon);

  // A store icon is keyed by preset, so a versioned member like mariadb-11-8
  // draws its family's mark. Preset and colour travel together: the mark is a
  // silhouette with no colours of its own, and the declared colour is what
  // gives it its tone.
  const resolvedPreset = $derived(preset ?? meta?.preset ?? name);
  const storeIcon = $derived($serviceIcons[resolvedPreset]);
  const brand = $derived(brandTintStyle(color ?? meta?.color));

  const tint = $derived(brand ? 'mark-tint' : tintFor(resolvedCategory));
  const box = $derived(compact ? 'w-8 h-8' : 'w-9 h-9 transition-transform group-hover:scale-105');
  const glyph = $derived(compact ? 'w-4 h-4' : 'w-5 h-5');
</script>

<span class="shrink-0 inline-flex items-center justify-center rounded-lg {tint} {box}" style={brand}>
  {#if storeIcon}
    <span class="mark-glyph {glyph}">{@html storeIcon}</span>
  {:else}
    <svg class={glyph} fill="none" stroke="currentColor" viewBox="0 0 24 24"
      >{@html dashboardIconSvg(name, resolvedIcon)}</svg
    >
  {/if}
</span>
