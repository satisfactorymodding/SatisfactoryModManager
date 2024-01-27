import type { ModalComponent } from '@skeletonlabs/skeleton';

import ProgressModal from './ProgressModal.svelte';
import ServerManager from './ServerManager.svelte';
import CacheLocationPicker from './settings/CacheLocationPicker.svelte';
import AddProfile from './profiles/AddProfile.svelte';
import ImportProfile from './profiles/ImportProfile.svelte';
import UpdatesModal from './updates/UpdatesModal.svelte';
import SMMUpdateReady from './smmUpdate/SMMUpdateReady.svelte';
import SMMUpdateDownload from './smmUpdate/SMMUpdateDownload.svelte';

// We can only store here modals (or modal instances) that do not require additional props
export const modalRegistry = {
  progress: { ref: ProgressModal } as ModalComponent,
  serverManager: { ref: ServerManager } as ModalComponent,
  cacheLocationPicker: { ref: CacheLocationPicker } as ModalComponent,
  addProfile: { ref: AddProfile } as ModalComponent,
  importProfile: { ref: ImportProfile } as ModalComponent,
  modUpdates: { ref: UpdatesModal } as ModalComponent,
  smmUpdateDownload: { ref: SMMUpdateDownload } as ModalComponent,
  smmUpdateReady: { ref: SMMUpdateReady } as ModalComponent,
};
							