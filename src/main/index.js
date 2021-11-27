import {
  app, BrowserWindow, ipcMain, shell, screen, session, dialog,
} from 'electron';
import path from 'path';
import http from 'http';
import { autoUpdater } from 'electron-updater';
import { Server } from 'socket.io';
import { getSetting, saveSetting } from '../settings';
import './differentialUpdateProgress';

require('electron-debug')({ isEnabled: true, showDevTools: false });

process.env.SMM_API_USERAGENT = process.env.NODE_ENV !== 'development' ? app.name : 'SMM-dev';
process.env.SMM_API_USERAGENT_VERSION = process.env.NODE_ENV !== 'development' ? app.getVersion() : 'development';

/**
 * Set `__static` path to static files in production
 * https://simulatedgreg.gitbooks.io/electron-vue/content/en/using-static-assets.html
 */
if (process.env.NODE_ENV !== 'development') {
  global.__static = path.join(__dirname, '/static').replace(/\\/g, '\\\\');
}

app.allowRendererProcessReuse = false;

/** @type { BrowserWindow } */
let mainWindow;
const mainURL = 'http://localhost:9080';
const mainFile = path.resolve(__dirname, 'index.html');

function sendToWindow(channel, ...args) {
  if (mainWindow && mainWindow.webContents) {
    mainWindow.webContents.send(channel, ...args);
  }
}

function openedByUrl(url) {
  if (url) {
    sendToWindow('openedByUrl', url);
  }
}

const normalSize = getSetting('normalSize', {
  width: 550,
  height: 850,
});
const minNormalSize = {
  width: 550,
  height: 650,
};
const expandedSize = getSetting('expandedSize', {
  width: 1575,
  height: 850,
});
const minExpandedSize = {
  width: 1225,
  height: 650,
};

let isExpanded = false;
let isChangingExpanded = false;

function updateSize() {
  const size = isExpanded ? expandedSize : normalSize;
  const minSize = isExpanded ? minExpandedSize : minNormalSize;
  mainWindow.setMinimumSize(minSize.width, minSize.height); // https://github.com/electron/electron/issues/15560#issuecomment-451395078
  mainWindow.setMaximumSize(isExpanded ? 2147483647 : normalSize.width, 2147483647);
  mainWindow.setSize(size.width, size.height, true);
}

function createWindow() {
  const frame = process.platform === 'linux';
  global.frame = frame;
  const windowLocation = getSetting('windowLocation', {});
  mainWindow = new BrowserWindow({
    x: windowLocation.x,
    y: windowLocation.y,
    width: normalSize.width,
    height: normalSize.height,
    minHeight: minNormalSize.height,
    minWidth: minNormalSize.width,
    maxWidth: normalSize.width,
    useContentSize: true,
    webPreferences: {
      nodeIntegration: true,
      contextIsolation: false,
      enableRemoteModule: true,
    },
    frame,
    show: false,
    icon: process.platform === 'linux' ? path.join(__dirname, '../../icons/64x64.png') : undefined, // https://github.com/AppImage/AppImageKit/wiki/Bundling-Electron-apps
  });

  mainWindow.webContents.on('did-finish-load', () => {
    mainWindow.show();
  });

  app.applicationMenu = null;

  if (getSetting('maximized', false)) {
    mainWindow.maximize();
  }

  if (process.env.NODE_ENV !== 'production') {
    mainWindow.loadURL(mainURL);
  } else {
    mainWindow.loadFile(mainFile);
  }

  mainWindow.on('resize', () => {
    if (!isChangingExpanded) {
      normalSize.height = mainWindow.getBounds().height;
      expandedSize.height = mainWindow.getBounds().height;
      if (isExpanded) {
        expandedSize.width = mainWindow.getBounds().width;
      }
      saveSetting('normalSize', normalSize);
      saveSetting('expandedSize', expandedSize);
    }
  });

  mainWindow.on('maximize', () => {
    saveSetting('maximized', true);
  });

  mainWindow.on('unmaximize', () => {
    saveSetting('maximized', false);
  });

  mainWindow.on('move', () => {
    saveSetting('windowLocation', { x: mainWindow.getBounds().x, y: mainWindow.getBounds().y });
  });

  ipcMain.on('openDevTools', () => {
    mainWindow.webContents.openDevTools();
  });

  mainWindow.on('closed', () => {
    mainWindow = null;
  });

  mainWindow.webContents.on('new-window', (event, url) => {
    event.preventDefault();
    shell.openExternal(url);
  });

  ipcMain.handle('saveDialog', (event, options) => dialog.showSaveDialogSync(mainWindow, options));

  ipcMain.handle('getVersion', () => app.getVersion());

  ipcMain.handle('hasFrame', () => frame);

  ipcMain.handle('minimize', () => mainWindow.minimize());
  ipcMain.handle('maximize', () => mainWindow.maximize());
  ipcMain.handle('unmaximize', () => mainWindow.unmaximize());
  ipcMain.handle('isMaximized', () => mainWindow.isMaximized());
  ipcMain.handle('close', () => mainWindow.close());
}

let isAutoUpdateTarget = true; // will be set to false if checkForUpdates errors
let isDownloadingUpdate = false;
let quitWaitingForUpdate = false;
let hasUpdate = false;

function isNetworkError(errorObject) {
  return errorObject.message === 'net::ERR_INTERNET_DISCONNECTED'
      || errorObject.message === 'net::ERR_PROXY_CONNECTION_FAILED'
      || errorObject.message === 'net::ERR_CONNECTION_RESET'
      || errorObject.message === 'net::ERR_CONNECTION_CLOSE'
      || errorObject.message === 'net::ERR_NAME_NOT_RESOLVED'
      || errorObject.message === 'net::ERR_CONNECTION_TIMED_OUT';
}

if (app.requestSingleInstanceLock()) {
  app.on('second-instance', (e, argv) => {
    if (process.platform === 'win32') {
      openedByUrl(argv.find((arg) => arg.startsWith('smmanager:')));
    }
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore();
      mainWindow.focus();
    }
  });

  app.on('ready', () => {
    session.defaultSession.webRequest.onBeforeRequest({ urls: ['https://www.youtube.com/get_video_info*'] }, (details, callback) => { // YT doesn't allow files to load embeds on purpose, even though it works!
      if (!details.url.includes('get_video_info') || !details.url.includes('&ancestor_origins=file%3A%2F%2F')) {
        callback({});
        return;
      }
      callback({ redirectURL: details.url.replace('&ancestor_origins=file%3A%2F%2F', '&ancestor_origins=http%3A%2F%2Flocalhost%3A9080').replace('&eurl', '&eurl=http%3A%2F%2Flocalhost%3A9080%2F') });
    });
    createWindow();
  });

  app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
      if (hasUpdate) {
        if (!isDownloadingUpdate) {
          autoUpdater.quitAndInstall(true, true);
        } else {
          quitWaitingForUpdate = true;
        }
      } else {
        app.quit();
      }
    }
  });

  app.on('activate', () => {
    if (mainWindow === null) {
      createWindow();
    }
  });

  ipcMain.once('vue-ready', () => {
    if (process.platform === 'win32') {
      openedByUrl(process.argv.find((arg) => arg.startsWith('smmanager:')));
    }
  });

  ipcMain.on('expand', () => {
    isChangingExpanded = true;
    isExpanded = true;
    updateSize();
    const windowScreen = screen.getDisplayMatching(mainWindow.getBounds());
    if (mainWindow.getBounds().x + mainWindow.getBounds().width > windowScreen.workArea.x + windowScreen.workArea.width) {
      mainWindow.setPosition(windowScreen.workArea.x + windowScreen.workArea.width - mainWindow.getBounds().width, mainWindow.getBounds().y, true);
    }
    isChangingExpanded = false;
  });

  ipcMain.on('unexpand', () => {
    isChangingExpanded = true;
    isExpanded = false;
    updateSize();
    const windowScreen = screen.getDisplayMatching(mainWindow.getBounds());
    if (mainWindow.getBounds().x + mainWindow.getBounds().width < windowScreen.workArea.x) {
      mainWindow.setPosition(windowScreen.workArea.x, mainWindow.getBounds().y, true);
    }
    isChangingExpanded = false;
  });

  autoUpdater.fullChangelog = true;

  autoUpdater.on('update-downloaded', () => {
    sendToWindow('updateDownloaded');
    if (quitWaitingForUpdate) {
      autoUpdater.quitAndInstall(true, false);
    } else {
      isDownloadingUpdate = false;
    }
  });

  autoUpdater.on('download-progress', (info) => {
    sendToWindow('updateDownloadProgress', info);
  });

  autoUpdater.on('error', () => {
    sendToWindow('updateNotAvailable');
    if (quitWaitingForUpdate) {
      app.quit();
    }
  });

  ipcMain.on('checkForUpdates', () => {
    if (isAutoUpdateTarget) {
      autoUpdater.checkForUpdates().catch((e) => {
        console.log(`Error checking for updates: ${e}`);
      });
    } else {
      sendToWindow('updateNotAvailable');
    }
  });

  autoUpdater.on('error', (_, err) => {
    sendToWindow('updateNotAvailable');
    isDownloadingUpdate = false;
    if (!err.includes('ENOENT') && !isNetworkError(err)) {
      console.error(err);
      // sendToWindow('autoUpdateError', err);
    } else {
      isAutoUpdateTarget = false;
    }
  });

  autoUpdater.on('update-not-available', () => {
    sendToWindow('updateNotAvailable');
  });

  autoUpdater.on('update-available', (updateInfo) => {
    sendToWindow('updateAvailable', updateInfo);
    isDownloadingUpdate = true;
    hasUpdate = true;
  });

  if (!app.isDefaultProtocolClient('smmanager')) {
    app.setAsDefaultProtocolClient('smmanager');
  }

  app.on('will-finish-launching', () => {
    app.on('open-url', (event, url) => {
      event.preventDefault();
      openedByUrl(url);
    });
  });

  const srv = http.createServer();
  srv.listen(33642, '127.0.0.1');
  const wss = new Server(srv, { path: '/' });
  wss.on('connection', (socket) => {
    socket.on('installedMods', () => {
      ipcMain.once('installedMods', (event, installedMods) => {
        const result = {};
        Object.entries(installedMods).forEach(([item, data]) => { result[item] = data.version; });
        socket.emit('installedMods', result);
      });
      mainWindow.webContents.send('installedMods');
    });
  });
} else {
  app.quit();
}
