<script lang="ts">
  import { GetAnnouncementsDocument } from '$lib/generated';
  import { getContextClient, queryStore } from '@urql/svelte';
  import Announcement from './Announcement.svelte';
  import Carousel from 'svelte-carousel';
  import { viewedAnnouncements } from '$lib/store/generalStore';
  import Tooltip, { Wrapper } from '@smui/tooltip';

  const client = getContextClient();

  const announcementsStore = queryStore({
    query: GetAnnouncementsDocument,
    client,
    requestPolicy: 'cache-and-network',
  });

  setInterval(() => {
    announcementsStore.pause();
    announcementsStore.resume();
  }, 1000 * 60 * 5);

  $: announcements = $announcementsStore.data?.getAnnouncements ?? [];

  let hovered = false;
  let currentIndex = 0;

  $: if(currentIndex < announcements.length) {
    setTimeout(() => {
      setAnnouncementViewed(announcements[currentIndex]);
    }, 6000);
  }

  $: if(hovered) {
    setAnnouncementViewed(announcements[currentIndex]);
  }

  function pageChange(e: CustomEvent<number>) {
    currentIndex = e.detail;
  }

  function setAnnouncementViewed(announcement: Pick<Announcement, 'id'>) {
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