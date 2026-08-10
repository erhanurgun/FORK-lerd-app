<script lang="ts">
  import { dashboardIconSvg } from '$lib/dashboardIcons';
  import { serviceMeta } from '$lib/serviceMeta';
  import { serviceIcons } from '$stores/serviceIcons';
  import { brandTintStyle } from '$lib/brandTint';
  import { tintFor, inkTintFor, asCategory } from '$lib/presetCategories';

  interface Props {
    name: string;
    category?: string;
    icon?: string;
    color?: string;
    preset?: string;
    compact?: boolean;
    // Drops the tinted plate and draws the mark on its own, larger, taking the
    // brand tone as its ink. For a header, where the plate is one box too many.
    bare?: boolean;
    // A mark in a strip of icons (the rail) answers to their hover and active
    // colours instead, so it carries no tone of its own.
    tint?: boolean;
  }
  let {
    name,
    category,
    icon,
    color,
    preset,
    compact = false,
    bare = false,
    tint = true
  }: Props = $props();

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
  const brand = $derived(tint ? brandTintStyle(color ?? meta?.color) : '');

  const tone = $derived(
    !tint
      ? ''
      : bare
        ? brand
          ? 'mark-brand'
          : inkTintFor(resolvedCategory)
        : brand
          ? 'mark-tint'
          : tintFor(resolvedCategory)
  );
  const box = $derived(
    bare ? '' : compact ? 'w-8 h-8' : 'w-9 h-9 transition-transform group-hover:scale-105'
  );
  const glyph = $derived(
    bare ? (compact ? 'w-5 h-5' : 'w-7 h-7') : compact ? 'w-4 h-4' : 'w-5 h-5'
  );
</script>

<span
  class="shrink-0 inline-flex items-center justify-center {bare ? '' : 'rounded-lg'} {tone} {box}"
  style={brand}
>
  {#if storeIcon}
    <span class="mark-glyph {glyph}">{@html storeIcon}</span>
  {:else}
    <svg class={glyph} fill="none" stroke="currentColor" viewBox="0 0 24 24"
      >{@html dashboardIconSvg(name, resolvedIcon)}</svg
    >
  {/if}
</span>
