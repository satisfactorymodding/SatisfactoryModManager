<script lang="ts">
  import { AnnouncementImportance, GetAnnouncementsDocument, SmrHealthcheckDocument, type Announcement as AnnouncementType } from '$lib/generated';
  import { getContextClient, queryStore } from '@urql/svelte';
  import Announcement from './Announcement.svelte';
  import Carousel from 'svelte-carousel';
  import { viewedAnnouncements } from '$lib/store/generalStore';
  import Tooltip, { Wrapper } from '@smui/tooltip';
  import { offline } from '$lib/store/settingsStore';

  const client = getContextClient();

  $: healthcheckStore = queryStore({
    query: SmrHealthcheckDocument,
    client,
    requestPolicy: 'network-only',
    pause: !!$offline,
  });

  let healthcheckGrace = true;
  $: if($offline !== null && !$offline) {
    setTimeout(() => {
      healthcheckGrace = false;
    }, 500);
  } else {
    healthcheckGrace = true;
  }

  $: healthcheck = $offline || !!$healthcheckStore.data?.getMods?.count || healthcheckGrace;

  const announcementsStore = queryStore({
    query: GetAnnouncementsDocument,
    client,
    pause: !!$offline,
    requestPolicy: 'cache-and-network',
  });

  setInterval(() => {
    announcementsStore.pause();
    announcementsStore.resume();
    healthcheckStore.pause();
    healthcheckStore.resume();
  }, 1000 * 60 * 5);

  const offlineAnnouncement: AnnouncementType = {
    id: '__offline__',
    message: 'You are currently offline. Some features may be unavailable.',
    importance: AnnouncementImportance.Info,
  };

  const healthcheckFailAnnouncement: AnnouncementType = {
    id: '__healthcheck__',
    message: 'Could not reach ficsit.app. Check your internet connection or consider using the offline mode.',
    importance: AnnouncementImportance.Warning,
  };

  $: customAnnouncements = [
    healthcheck ? null : healthcheckFailAnnouncement,
  ].filter((a) => a !== null) as AnnouncementType[];

  $: smrAnnouncements = $announcementsStore.data?.getAnnouncements ?? [];

  $: announcements = $offline === null ? [] : ($offline ? [offlineAnnouncement] : [...smrAnnouncements, ...customAnnouncements]);

  let hovered = false;
  let currentIndex = 0;

  $: if(currentIndex < announcements.length) {
    const viewedAnnouncement = announcements[currentIndex];
    setTimeout(() => {
      setAnnouncementViewed(viewedAnnouncement);
    }, 6000);
  }

  $: if(hovered) {
    setAnnouncementViewed(announcements[currentIndex]);
  }

  function pageChange(e: CustomEvent<number>) {
    currentIndex = e.detail;
  }

  function setAnnouncementViewed(announcement: Pick<Announcement, 'id'>) {
    if(announcement.id === offlineAnnouncement.id) return;
    if(announcement.id === healthcheckFailAnnouncement.id) return;
    if(!$viewedAnnouncements.some((id) => announcement.id === id)) {
      viewedAnnouncements.update((ids) => [...ids, announcement.id]);
    }
  }
</script>

{#if announcements.length > 0}
  <div class="w-full" on:mouseenter={() => hovered = true} on:mouseleave={() => hovered = false}>
    <Carousel
      autoplayDuration={hovered ? 1e100 : 5000}
      duration={300}
      autoplay={!hovered && announcements.length > 1}
      dots={false}
      arrows={false}
      swiping={false}
      on:pageChange={pageChange}
    >
      {#each announcements as announcement}
        <div class="w-full shrink-0">
          <Wrapper>
            <Announcement {announcement} />
            <Tooltip surface$class="max-w-xl text-base">
              <span class="text-xl">{announcement.message}</span>
            </Tooltip>
          </Wrapper>
        </div>
      {/each}
    </Carousel>
  </div>
{/if}