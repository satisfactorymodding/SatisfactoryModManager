<template>
  <v-row
    align="center"
    no-gutters
    class="py-1 announcement-bg"
    :class="{ [`announcement-${importanceLower}`]: true,
              'announcement-new': isNew } "
  >
    <v-row
      align="center"
      no-gutters
      class="mx-1 announcement-bg-text"
      :class="`announcement-${importanceLower}`"
    >
      <v-col
        cols="auto"
        class="pl-4 pr-4"
      >
        <v-icon
          style="font-size: 32px !important"
          :color="announcementFgColor"
        >
          {{ announcementIcon }}
        </v-icon>
      </v-col>
      <v-col
        cols="10"
        :style="`color: ${announcementFgColor} !important;`"
      >
        <!-- eslint-disable-next-line vue/no-v-html -->
        <span v-html="announcementMessageHTML" />
      </v-col>
    </v-row>
  </v-row>
</template>

<script>
import { markdownAsElement } from '@/utils';
import { getSetting, saveSetting } from '~/settings';

export default {
  props: {
    announcement: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isNew: false,
    };
  },
  computed: {
    importanceLower() {
      return this.announcement.importance.toLowerCase();
    },
    announcementIcon() {
      switch (this.importanceLower) {
        case 'alert':
          return 'mdi-alert-outline';
        case 'warning':
          return 'mdi-alert-outline';
        case 'fix':
          return 'mdi-information-outline';
        case 'info':
        default:
          return 'mdi-information-outline';
      }
    },
    announcementBgColor() {
      switch (this.importanceLower) {
        case 'alert':
          return '#bb2828';
        case 'warning':
          return '#ffcc00';
        case 'fix':
          return '#249a20';
        case 'info':
        default:
          return '#5daec5';
      }
    },
    announcementFgColor() {
      switch (this.importanceLower) {
        default:
          return 'white';
      }
    },
    announcementMessageElement() {
      return markdownAsElement(this.announcement.message);
    },
    announcementMessageHTML() {
      const el = this.announcementMessageElement;
      const links = el.getElementsByTagName('a');
      for (let i = 0; i < links.length; i += 1) {
        links[i].target = '_blank';
        links[i].style += ';color: unset;';
      }
      // get content from inside paragraph
      return el.children[1].children[0].innerHTML;
    },
  },
  mounted() {
    this.isNew = true || !getSetting('viewedAnnouncements', []).includes(this.announcement.id);
    if (this.isNew) {
      setTimeout(() => {
        const viewedAnnouncements = getSetting('viewedAnnouncements', []);
        viewedAnnouncements.push(this.announcement.id);
        saveSetting('viewedAnnouncements', viewedAnnouncements);
        this.isNew = false;
      }, 60 * 1000);
    }
  },
};
</script>

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
  animation: slide 6s linear infinite;
  will-change: background-position;
}

.announcement-new.announcement-bg {
  background-image: repeating-linear-gradient(
    45deg,
    transparent,
    transparent 20px,
    var(--colour2) 20px,
    var(--colour2) 40px
  );
}
</style>
