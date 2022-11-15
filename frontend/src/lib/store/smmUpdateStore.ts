import { readable } from 'svelte/store';
import { compare } from 'semver';
import { EventsOn } from '$wailsjs/runtime/runtime';


export interface SMMUpdate {
  newVersion: string;
  changelogs: {
    version: string;
    changelog: string;
  }[];
}

export interface SMMUpdateProgress {
  downloaded: number;
  total: number;
  speed: number;
}
  

export const smmUpdate = readable<SMMUpdate | null>(null, (set) => {
  EventsOn('updateAvailable', (newVersion: string, changelogs: Record<string, string>) => {
    set({ newVersion, changelogs: Object.entries(changelogs).map(([version, changelog]) => ({ version, changelog })).sort((a, b) => -compare(a.version, b.version)) });
  });
});

export const smmUpdateProgress = readable<SMMUpdateProgress | null>(null, (set) => {
  const pastDownloaded: { downloaded: number; time: number } = [];
  const speedTimeframe = 1000 * 5;
  EventsOn('updateDownloadProgress', (downloaded: number, total: number) => {
    pastDownloaded.push({
      downloaded,
      time: Date.now(),
    });
    setTimeout(() => {
      pastDownloaded.shift();
    }, speedTimeframe);
    set({ downloaded, total, speed: pastDownloaded.length > 0 ? (downloaded - pastDownloaded[0].downloaded) / ((Date.now() - pastDownloaded[0].time) / 1000) : 0 });
  });
});

export const smmUpdateReady = readable<boolean>(false, (set) => {
  EventsOn('updateReady', () => {
    set(true);
  });
  EventsOn('updateAvailable', () => {
    set(false);
  });
});