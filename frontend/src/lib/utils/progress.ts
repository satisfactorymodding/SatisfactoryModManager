import { type Readable, get, writable } from 'svelte/store';

import { timeSeries } from './timeSeries';

import type { utils } from '$wailsjs/go/models';

export function progressStats(progress: Readable<utils.Progress | null>, options?: { statsInterval?: number; updateInterval?: { speed?: number; eta?: number } }): Readable<{ speed: number; eta: number | undefined }> {
  const finalOptions = {
    statsInterval: 5000,
    ...options,
    updateInterval: {
      speed: 0,
      eta: 0,
      ...options?.updateInterval,
    },
  };

  const series = timeSeries(finalOptions.statsInterval);
  
  const stats = writable({ speed: 0, eta: 0 as number | undefined });
  const lastStatsUpdate = { speed: 0, eta: 0 };
  
  progress.subscribe(($progress) => {
    if (!$progress) {
      series.clear();
    } else {    
      series.addValue($progress.current);
  
      const speed = series.getDerivative() ?? 0;
      const eta = speed !== 0 ? ($progress.total - $progress.current) / speed : undefined;
  
      if (Date.now() - lastStatsUpdate.speed > finalOptions.updateInterval.speed) {
        stats.set({ speed: speed, eta: get(stats).eta });
        lastStatsUpdate.speed = Date.now();
      }
      if (Date.now() - lastStatsUpdate.eta > finalOptions.updateInterval.eta) {
        stats.set({ speed: get(stats).speed, eta: eta });
        lastStatsUpdate.eta = Date.now();
      }
    }
  });

  return stats;
}
