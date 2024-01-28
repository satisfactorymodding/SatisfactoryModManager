<script lang="ts">
  import { getContextClient, queryStore } from '@urql/svelte';
  import Carousel from 'svelte-carousel';
  import { popup, type PopupSettings } from '@skeletonlabs/skeleton';

  import Tooltip from '$lib/components/Tooltip.svelte';
  import Announcement from '$lib/components/announcements/Announcement.svelte';
  import { viewedAnnouncements , offline } from '$lib/store/settingsStore';
  import { AnnouncementImportance, GetAnnouncementsDocument, SmrHealthcheckDocument, type Announcement as AnnouncementType } from '$lib/generated';
  import { SetAnnouncementViewed } from '$wailsjs/go/bindings/Settings';
  
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

  interface ViewableAnnouncement extends AnnouncementType {
    viewable: boolean;
  }

  const CUSTOM_ANNOUNCEMENT_REGEX = /^__.*__$/;

  const offlineAnnouncement: ViewableAnnouncement = {
    id: '__offline__',
    message: 'You are currently offline. Some features may be unavailable. (To reconnect, use Mod Manager Settings > Go Online)',
    importance: AnnouncementImportance.Info,
    viewable: false,
  };

  const healthcheckFailAnnouncement: ViewableAnnouncement = {
    id: '__healthcheck__',
    message: 'Could not reach ficsit.app. Check your internet connection or consider using the offline mode. (Mod Manager Settings > Go Offline)',
    importance: AnnouncementImportance.Warning,
    viewable: false,
  };

  $: customAnnouncements = [
    healthcheck ? null : healthcheckFailAnnouncement,
  ].filter((a) => a !== null) as ViewableAnnouncement[];

  $: smrAnnouncements = $announcementsStore.data?.getAnnouncements?.map((a) => ({ ...a, viewable: true })) ?? [];

  $: announcements = $offline === null ? [] : [...customAnnouncements, ...smrAnnouncements];

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

  function setAnnouncementViewed(announcement?: Pick<ViewableAnnouncement, 'id' | 'viewable'>) {
    if(!announcement?.viewable) return;
    if(!$viewedAnnouncements.some((id) => announcement.id === id)) {
      SetAnnouncementViewed(announcement.id);
    }
  }

  function goOnline() {
    $offline = false;
  }

  function goOffline() {
    $offline = true;
  }

  $: popupId = 'announcement';

  $: popupHover = {
    event: 'hover',
    target: popupId,
    middleware: {
      offset: 4,
    },
    placement: 'bottom',
  } satisfies PopupSettings;
</script>

<!-- the if gets executed before this is added to the DOM for some reason if this is below the ifs, so the use:popup would not find the element -->
<Tooltip disabled={!$offline && !announcements[currentIndex]} {popupId}>
  <!--
    fixed allows the popup to be displayed outside the bounds of the parent
    block opacity-0 ensure that the popup gets a width on first render, otherwise it will be 0px wide and floating-ui will "animate" widening it
    inert is required when popup is invisible to ignore mouse inputs
  -->
  <span>
    {#if $offline}
      {offlineAnnouncement.message}
    {:else}
      {announcements[currentIndex]?.message}
    {/if}
  </span>
</Tooltip>

{#if $offline}
  <div class="w-full" use:popup={popupHover}>
    <Announcement announcement={offlineAnnouncement}>
      <div class="flex pr-2">
        <span>{offlineAnnouncement.message}</span>
        <span
          class="text-yellow-400 font-bold underline cursor-pointer ml-auto"
          role="button"
          tabindex="0"
          on:click={goOnline}
          on:keypress={goOnline}>Go Online</span>
      </div>
    </Announcement>
  </div>
{:else if announcements.length > 0}
  <div class="w-full" role="alert" use:popup={popupHover}>
    <Carousel
      arrows={false}
      autoplay={!hovered && announcements.length > 1}
      autoplayDuration={hovered ? 1e100 : 5000}
      dots={false}
      duration={300}
      swiping={false}
      on:pageChange={pageChange}
    >
      {#each announcements as announcement}
        <div class="w-full shrink-0">
          {#if CUSTOM_ANNOUNCEMENT_REGEX.test(announcement.id)}
            <Announcement {announcement}>
              {#if announcement.id === '__healthcheck__'}
                <div class="flex pr-2">
                  <span>{announcement.message}</span>
                  <span
                    class="text-yellow-400 font-bold underline cursor-pointer ml-auto"
                    role="button"
                    tabindex="0"
                    on:click={goOffline}
                    on:keypress={goOffline}>Go Offline</span>
                </div>
              {:else}
                {announcement.message}
              {/if}
            </Announcement>
          {:else}
            <Announcement {announcement} />
          {/if}
        </div>
      {/each}
    </Carousel>
  </div>
{/if}
