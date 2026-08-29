import { writable } from 'svelte/store';
import { apiJson, apiFetch } from '$lib/api';

export const autostartEnabled = writable<boolean>(false);
export const startOnDashboardOpen = writable<boolean>(false);

interface SettingsResponse {
  autostart_on_login?: boolean;
  start_on_dashboard_open?: boolean;
}

export async function loadAutostart() {
  try {
    const res = await apiJson<SettingsResponse>('/api/settings');
    autostartEnabled.set(Boolean(res.autostart_on_login));
    startOnDashboardOpen.set(Boolean(res.start_on_dashboard_open));
  } catch {
    /* keep previous */
  }
}

export async function toggleAutostart(enable: boolean): Promise<boolean> {
  try {
    const res = await apiFetch('/api/settings/autostart', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: enable })
    });
    if (res.ok) autostartEnabled.set(enable);
    return res.ok;
  } catch {
    return false;
  }
}

export async function toggleStartOnDashboardOpen(enable: boolean): Promise<boolean> {
  try {
    const res = await apiFetch('/api/settings/start-on-open', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ enabled: enable })
    });
    if (res.ok) startOnDashboardOpen.set(enable);
    return res.ok;
  } catch {
    return false;
  }
}
