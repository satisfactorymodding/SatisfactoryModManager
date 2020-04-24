
import {
  app, BrowserWindow, ipcMain, shell, Menu, screen,
} from 'electron';
import path from 'path';
import { autoUpdater } from 'electron-updater';
import WebSocket from 'ws';


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

const isMac = process.platform === 'darwin';

const template = [
  // { role: 'appMenu' }
  ...(isMac ? [{
    label: app.name,
    submenu: [
      { role: 'about' },
      { type: 'separator' },
      { role: 'services' },
      { type: 'separator' },
      { role: 'hide' },
      { role: 'hideothers' },
      { role: 'unhide' },
      { type: 'separator' },
      { role: 'quit' },
    ],
  }] : []),
  // { role: 'fileMenu' }
  {
    label: 'File',
    submenu: [
      isMac ? { role: 'close' } : { role: 'quit' },
    ],
  },
  // { role: 'viewMenu' }
  {
    label: 'View',
    submenu: [
      { role: 'reload' },
      { role: 'forcereload' },
      { role: 'toggledevtools' },
      { type: 'separator' },
      { role: 'resetzoom' },
      { role: 'zoomin' },
      { role: 'zoomout' },
      { type: 'separator' },
      { role: 'togglefullscreen' },
    ],
  },
  {
    role: 'help',
    submenu: [
      {
        label: `SMLauncher v${app.getVersion()}`,
        enabled: false,
      },
      {
        label: 'Join the Satisfactory Modding Discord',
        click: () => {
          shell.openExternal('https://discord.gg/TShj39G');
        },
      },
      {
        label: 'Toggle Debug Mode',
        click: () => {
          mainWindow.webContents.send('toggleDebug');
        },
      },
      {
        label: 'Clear Cache',
        click: () => {
          mainWindow.webContents.send('clearCache');
        },
      },
    ],
  },
];

const menu = Menu.buildFromTemplate(template);
Menu.setApplicationMenu(menu);

const initialWidth = 500;
const initialHeight = 700;

function createWindow() {
  /**
   * Initial window options
   */
  mainWindow = new BrowserWindow({
    useContentSize: true,
    height: initialHeight,
    width: initialWidth,
    minWidth: initialWidth,
    minHeight: initialHeight,
    maxWidth: initialWidth,
    webPreferences: {
      nodeIntegration: true,
    },
    frame: false,
    show: false,
  });

  ipcMain.once('vue-ready', () => {
    if (process.platform === 'win32') {
      openedByUrl(process.argv.find((arg) => arg.startsWith('smlauncher:')));
    }
  });

  if (process.env.NODE_ENV === 'development') {
    mainWindow.loadURL(mainURL);
  } else {
    mainWindow.loadFile(mainFile);
  }

  ipcMain.on('openDevTools', () => {
    mainWindow.webContents.openDevTools();
  });

  mainWindow.once('ready-to-show', () => {
    mainWindow.show();
  });

  ipcMain.on('open-side-panel', () => {
    const widthScreen = screen.getPrimaryDisplay().workAreaSize.width;
    const heightScreen = screen.getPrimaryDisplay().workAreaSize.height;
    mainWindow.setMinimumSize(initialWidth + 500, initialHeight);
    mainWindow.setMaximumSize(widthScreen, heightScreen);
    mainWindow.setBounds(
      {
        x: 10, y: 10, width: widthScreen - 20, height: heightScreen - 20,
      },
    );
  });

  ipcMain.on('close-side-panel', () => {
    const widthScreen = screen.getPrimaryDisplay().workAreaSize.width;
    const heightScreen = screen.getPrimaryDisplay().workAreaSize.height;
    const MaxHeight = heightScreen - 20;
    mainWindow.setMinimumSize(initialWidth, initialHeight);
    mainWindow.setMaximumSize(initialWidth, heightScreen);
    mainWindow.setBounds(
      {
        x: (widthScreen / 2) - (initialWidth / 2),
        y: mainWindow.getPosition()[1],
        width: initialWidth,
        height: MaxHeight,
      },
    );
  });

  ipcMain.on('display-app-menu', () => {
    if (mainWindow) {
      menu.popup({
        window: mainWindow,
        x: 10,
        y: 10,
      });
    }
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

if (app.requestSingleInstanceLock()) {
  app.on('second-instance', (e, argv) => {
    if (process.platform === 'win32') {
      openedByUrl(argv.find((arg) => arg.startsWith('smlauncher:')));
    }
    if (mainWindow) {
      if (mainWindow.isMinimized()) mainWindow.restore();
      mainWindow.focus();
    }
  });

  app.on('ready', createWindow);

  app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
      if (!isDownloadingUpdate) {
        app.quit();
      } else {
        quitWaitingForUpdate = true;
      }
    }
  });

  app.on('activate', () => {
    if (mainWindow === null) {
      createWindow();
    }
  });

  autoUpdater.on('update-downloaded', () => {
    if (quitWaitingForUpdate) {
      autoUpdater.quitAndInstall(true);
    } else {
      isDownloadingUpdate = false;
    }
  });

  autoUpdater.on('error', () => {
    if (quitWaitingForUpdate) {
      app.quit();
    }
  });

  autoUpdater.on('update-available', (updateInfo) => {
    mainWindow.webContents.send('update-available', updateInfo);
    isDownloadingUpdate = true;
  });

  app.on('ready', () => {
    if (process.env.NODE_ENV === 'production') {
      autoUpdater.checkForUpdates();
    }
  });

  if (!app.isDefaultProtocolClient('smlauncher')) {
    app.setAsDefaultProtocolClient('smlauncher');
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
