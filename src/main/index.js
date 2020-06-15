
import {
  app, BrowserWindow, ipcMain, shell,
} from 'electron';
import path from 'path';
import { autoUpdater } from 'electron-updater';
import WebSocket from 'ws';


process.env.SMM_API_USERAGENT = process.env.NODE_ENV !== 'development' ? app.name : 'SMM-dev';
process.env.SMM_API_USERAGENT_VERSION = process.env.NODE_ENV !== 'development' ? app.getVersion() : 'development';

/**
 * Set `__static` path to static files in production
 * https://simulatedgreg.gitbooks.io/electron-vue/content/en/using-static-assets.html
 */
if (process.env.NODE_ENV !== 'development') {
  global.__static = path.join(__dirname, '/static').replace(/\\/g, '\\\\');
}

let mainWindow;
const mainURL = 'http://localhost:9080';
const mainFile = path.resolve(__dirname, 'index.html');

function openedByUrl(url) {
  if (url) {
    mainWindow.webContents.send('openedByUrl', url);
  }
}

const normalSize = {
  width: 500,
  height: 858,
};
const expandedSize = {
  width: 1575,
  height: 858,
};

let isExpanded = false;

function updateSize() {
  const size = isExpanded ? expandedSize : normalSize;
  mainWindow.setMinimumSize(size.width, size.height); // https://github.com/electron/electron/issues/15560#issuecomment-451395078
  mainWindow.setSize(size.width, size.height, true);
}

app.commandLine.appendSwitch('high-dpi-support', 1);
app.commandLine.appendSwitch('force-device-scale-factor', 1);

function createWindow() {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    width: normalSize.width,
    height: normalSize.height,
    useContentSize: true,
    minHeight: normalSize.height,
    minWidth: normalSize.width,
    webPreferences: {
      nodeIntegration: true,
    },
    frame: false,
    resizable: false,
  });

  if (process.env.NODE_ENV === 'development') {
    mainWindow.loadURL(mainURL);
  } else {
    mainWindow.loadFile(mainFile);
  }

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
}

let isDownloadingUpdate = false;
let quitWaitingForUpdate = false;
let hasUpdate = false;

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

  app.on('ready', createWindow);

  app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
      if (hasUpdate) {
        if (!isDownloadingUpdate) {
          autoUpdater.quitAndInstall(false);
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
    isExpanded = true;
    updateSize();
  });

  ipcMain.on('unexpand', () => {
    isExpanded = false;
    updateSize();
  });

  autoUpdater.fullChangelog = true;

  autoUpdater.on('update-downloaded', () => {
    mainWindow.webContents.send('updateDownloaded');
    if (quitWaitingForUpdate) {
      autoUpdater.quitAndInstall(false);
    } else {
      isDownloadingUpdate = false;
    }
  });

  autoUpdater.on('download-progress', (info) => {
    mainWindow.webContents.send('updateDownloadProgress', info);
  });

  autoUpdater.on('error', () => {
    mainWindow.webContents.send('updateNotAvailable');
    if (quitWaitingForUpdate) {
      app.quit();
    }
  });

  ipcMain.on('checkForUpdates', () => {
    autoUpdater.checkForUpdates();
  });

  autoUpdater.on('update-not-available', () => {
    mainWindow.webContents.send('updateNotAvailable');
  });

  autoUpdater.on('update-available', (updateInfo) => {
    mainWindow.webContents.send('updateAvailable', updateInfo);
    isDownloadingUpdate = true;
    hasUpdate = true;
  });

  app.on('ready', () => {
    if (process.env.NODE_ENV === 'production') {
      autoUpdater.checkForUpdates();
    }
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

  const wss = new WebSocket.Server({ port: 33642 });
  wss.on('connection', (ws) => {
    ws.on('message', (message) => {
      console.log('received: %s', message);
    });
  });
} else {
  app.quit();
}
