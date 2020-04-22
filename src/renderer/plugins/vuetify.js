import Vue from 'vue';
import Vuetify from 'vuetify/lib';
import colors from 'vuetify/lib/util/colors';
import {
  mdiCog,
  mdiWindowMinimize,
  mdiWindowClose,
  mdiTrashCan,
  mdiPlus,
  mdiEye,
  mdiEyeOff,
  mdiStar,
  mdiStarOutline,
  mdiImport,
  mdiGamepadVariant,
  mdiSort,
  mdiCogs,
  mdiFilter,
  mdiSortAscending,
  mdiSortDescending,
  mdiDownload,
  mdiBellCancel,
  mdiBookmark,
  mdiBookmarkCheck,
  mdiKeyboardBackspace,
} from '@mdi/js';
import '../assets/sass/main.scss';

Vue.use(Vuetify);

const opts = {
  theme: {
    options: {
      customProperties: true,
    },
    dark: true,
    themes: {
      light: {
        anchor: '#000000',
        'color-1': colors.grey.lighten5,
        'color-2': colors.grey.lighten4,
        'color-3': colors.red.lighten4,
      },
      dark: {
        anchor: '#ffffff',
        'color-1': colors.grey.darken4,
        'color-2': colors.grey.darken3,
        'color-3': colors.red.darken4,
      },
    },
  },
  icons: {
    iconfont: 'mdiSvg',
    values: {
      settingIcon: mdiCog,
      minimizeIcon: mdiWindowMinimize,
      closeIcon: mdiWindowClose,
      deleteIcon: mdiTrashCan,
      addIcon: mdiPlus,
      visibleIcon: mdiEye,
      visibleOffIcon: mdiEyeOff,
      starOnIcon: mdiStar,
      starOffIcon: mdiStarOutline,
      showMoreIcon: mdiImport,
      gamepadIcon: mdiGamepadVariant,
      sortIcon: mdiSort,
      filterIcon: mdiFilter,
      configIcon: mdiCogs,
      sortAscIcon: mdiSortAscending,
      sortDescIcon: mdiSortDescending,
      updateIcon: mdiDownload,
      ignoreIcon: mdiBellCancel,
      compatibleIcon: mdiBookmarkCheck,
      notCompatibleIcon: mdiBookmark,
      backIcon: mdiKeyboardBackspace,
    },
  },
};

export default new Vuetify(opts);
