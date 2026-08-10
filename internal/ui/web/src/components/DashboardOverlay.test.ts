import { render, screen, waitFor } from '@testing-library/svelte';
import { describe, it, expect, beforeEach } from 'vitest';
import DashboardOverlay from './DashboardOverlay.svelte';
import { dashboardOpen, openDocs } from '../stores/dashboard';
import { profilerEnabled } from '../stores/profiler';

function openProfiler() {
  dashboardOpen.set({
    name: 'profiler',
    label: 'Profiler',
    dashboard: '/_spx/?SPX_UI_URI=/'
  });
}

describe('DashboardOverlay', () => {
  beforeEach(() => {
    dashboardOpen.set(null);
    profilerEnabled.set(false);
  });

  it('disables Back until the embedded iframe has somewhere to go back to', () => {
    openProfiler();
    render(DashboardOverlay);

    // Freshly opened: the SPX iframe has no internal history yet, so Back is a
    // dead end. It must be disabled rather than silently tear down the overlay.
    expect(screen.getByTitle('Back')).toBeDisabled();
  });

  it('shows the profiler toggle as off: muted, not pressed, no live dot', () => {
    profilerEnabled.set(false);
    openProfiler();
    const { container } = render(DashboardOverlay);

    const btn = screen.getByRole('button', { name: /start profiling/i });
    expect(btn.getAttribute('aria-pressed')).toBe('false');
    expect(btn.className).not.toMatch(/emerald/);
    expect(container.querySelector('.animate-pulse')).toBeNull();
  });

  it('shows the profiler toggle as on: emerald, pressed, live pulsing dot', () => {
    profilerEnabled.set(true);
    openProfiler();
    const { container } = render(DashboardOverlay);

    const btn = screen.getByRole('button', { name: /stop profiling/i });
    expect(btn.getAttribute('aria-pressed')).toBe('true');
    expect(btn.className).toMatch(/emerald/);
    expect(container.querySelector('.animate-pulse')).not.toBeNull();
  });

  it('renders the documentation in place of an iframe', async () => {
    globalThis.fetch = (async () =>
      new Response(JSON.stringify({ pages: [] }), { status: 200 })) as unknown as typeof fetch;
    openDocs();
    const { container } = render(DashboardOverlay);

    expect(container.querySelector('iframe')).toBeNull();
    expect(screen.getByRole('searchbox')).toBeInTheDocument();
    // The header keeps a way out to the same page on the website.
    await waitFor(() =>
      expect(screen.getByTitle('Open in new tab').getAttribute('href')).toBe(
        'https://lerd.sh/getting-started/requirements'
      )
    );
  });

  it('collapses the SPX Configuration form by default on the control panel page', () => {
    openProfiler();
    render(DashboardOverlay);

    // The form starts hidden, so the header offers to show it.
    expect(screen.getByRole('button', { name: /show configuration/i })).toBeTruthy();
  });
});
