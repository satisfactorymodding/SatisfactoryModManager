import electronLogger from 'electron-log';
import { autoUpdater } from 'electron-updater';

autoUpdater.logger = electronLogger;
autoUpdater.logger.transports.file.level = 'info';

let diffDown = {
  percent: 0,
  bytesPerSecond: 0,
  total: 0,
  transferred: 0,
};
let diffDownHelper = {
  startTime: 0,
  lastTime: 0,
  lastSize: 0,
};

electronLogger.hooks.push((msg, transport) => {
  if (transport !== electronLogger.transports.console) {
    return msg;
  }

  let match = /Full: ([\d,.]+) ([GMKB]+), To download: ([\d,.]+) ([GMKB]+)/.exec(
    msg.data[0],
  );
  if (match) {
    let multiplier = 1;
    if (match[4] === 'KB') multiplier *= 1024;
    if (match[4] === 'MB') multiplier *= 1024 * 1024;
    if (match[4] === 'GB') multiplier *= 1024 * 1024 * 1024;

    diffDown = {
      percent: 0,
      bytesPerSecond: 0,
      total: Number(match[3].split(',').join('')) * multiplier,
      transferred: 0,
    };
    diffDownHelper = {
      startTime: Date.now(),
      lastTime: Date.now(),
      lastSize: 0,
    };
    return msg;
  }

  match = /download range: bytes=(\d+)-(\d+)/.exec(msg.data[0]);
  if (match) {
    const currentSize = Number(match[2]) - Number(match[1]);
    const currentTime = Date.now();
    const deltaTime = currentTime - diffDownHelper.startTime;

    diffDown.transferred += diffDownHelper.lastSize;
    diffDown.bytesPerSecond = Math.floor(
      (diffDown.transferred * 1000) / deltaTime,
    );
    diffDown.percent = (diffDown.transferred * 100) / diffDown.total;

    diffDownHelper.lastSize = currentSize;
    diffDownHelper.lastTime = currentTime;
    autoUpdater.emit('download-progress', diffDown);
    return msg;
  }
  return msg;
});
