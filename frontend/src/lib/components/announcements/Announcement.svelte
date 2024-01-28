<script lang="ts">
  import { mdiAlertOutline, mdiInformationOutline } from '@mdi/js';

  import { AnnouncementImportance, type Announcement } from '$lib/generated';
  import SvgIcon from '$lib/components/SVGIcon.svelte';
  import { viewedAnnouncements } from '$lib/store/settingsStore';

  export let announcement: Pick<Announcement, 'id' | 'importance' | 'message'>;

  $: importanceLower = announcement.importance.toLowerCase();
  $: isNew = !$viewedAnnouncements.includes(announcement.id);
  $: icon = (() => {
    switch (announcement.importance) {
      case AnnouncementImportance.Alert:
      case AnnouncementImportance.Warning:
        return mdiAlertOutline;
      case AnnouncementImportance.Fix:
      case AnnouncementImportance.Info:
      default:
        return mdiInformationOutline;
    }
  })();
</script>

<div class="announcement-{importanceLower} announcement-bg p-1.5 h-full" class:announcement-new={isNew}>
  <div class="flex items-center announcement-bg-text p-1 h-full">
    <SvgIcon class="w-8 h-8 mr-3 shrink-0" icon={icon} />
    <div class="grow wrap text-lg">
      <slot>
        {announcement.message}
      </slot>
    </div>
  </div>
</div>


<style scoped>
  .announcement-alert {
    --deg: 15deg;
  }
  .announcement-warning {
    --deg: 30deg;
  }
  .announcement-fix {
    --deg: 80deg;
  }
  .announcement-info {
    --deg: 186deg;
  }
  @keyframes slide {
    from {
      background-position-x: -113px;
    }
    to {
      background-position-x: 0px;
    }
  }
  .announcement-bg-text {
    --colour1: hsl(var(--deg) 100% 36%);
    background-color: var(--colour1);
    border-radius: 5px;
  }
  .announcement-bg {
    --colour1: hsl(var(--deg) 100% 36%);
    --colour2: hsl(var(--deg) 77% 24%);
    background-size: 200% 100% !important;
    background-color: var(--colour1);
    will-change: background-position;
    background-image: repeating-linear-gradient(
        45deg,
        transparent,
        transparent 20px,
        var(--colour2) 20px,
        var(--colour2) 40px
    );
  }
  .announcement-new.announcement-bg {
    animation: slide 6s linear infinite;
  }
  </style>
