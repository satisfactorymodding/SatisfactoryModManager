// eslint-disable-next-line max-classes-per-file
import { LogLevel, addLogger } from 'satisfactory-mod-manager-api';
import { getCacheFolder } from 'platform-folders';
import fs from 'fs';
import path from 'path';

export const logsDir = path.join(getCacheFolder(), 'SatisfactoryModManager', 'logs');
if (!fs.existsSync(logsDir)) {
  fs.mkdirSync(logsDir, { recursive: true });
}

class ConsoleLogger {
  constructor(minLevel) {
    this.minLevel = minLevel || LogLevel.INFO;
  }

  write(level, message) {
    if (level >= this.minLevel) {
      switch (level) {
        case LogLevel.DEBUG:
          console.log(message);
          break;
        case LogLevel.WARN:
          console.warn(message);
          break;
        case LogLevel.ERROR:
          console.error(message);
          break;
        case LogLevel.INFO:
        default:
          console.info(message);
          break;
      }
    }
  }
}

function formatDate(date) {
  return `${date.getFullYear().toString().padStart(4, '0')}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`;
}

class RollingFileLogger {
  constructor(dir, fileNameFormat, minLevel) {
    this.dir = dir;
    this.fileNameFormat = fileNameFormat;
    this.minLevel = minLevel || LogLevel.DEBUG;
    this.logFileWriter = fs.createWriteStream(this.getLogFilePath(), { flags: 'a', encoding: 'utf8', autoClose: true });
  }

  static formatLogFileName(fileName) {
    return fileName.replace('%DATE%', formatDate(new Date()));
  }

  getLogFilePath() {
    return path.join(this.dir, RollingFileLogger.formatLogFileName(this.fileNameFormat));
  }

  checkRoll() {
    if (this.logFileWriter.path !== this.getLogFilePath()) {
      this.logFileWriter.end('\n');
      this.logFileWriter = fs.createWriteStream(this.getLogFilePath(), { flags: 'a', encoding: 'utf8', autoClose: true });
      this.logFileWriter.write('\n');
    }
  }

  write(level, message) {
    if (level >= this.minLevel) {
      this.checkRoll();
      if (this.logFileWriter && this.logFileWriter.writable) {
        this.logFileWriter.write(message);
        this.logFileWriter.write('\n');
      }
    }
  }
}

export const consoleLogger = new ConsoleLogger();
export const fileLogger = new RollingFileLogger(logsDir, 'SatisfactoryModManager-%DATE%.log');

addLogger(consoleLogger);
addLogger(fileLogger);

function formatMessage(message) {
  if (message instanceof Error) {
    return `${message.message}\nTrace\n${message.stack}`;
  }
  if (typeof message === 'string') {
    return message;
  }
  return JSON.stringify(message);
}

function formatDateTime(date) {
  return `${date.getFullYear().toString().padStart(4, '0')}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}:${date.getSeconds().toString().padStart(2, '0')}:${date.getMilliseconds().toString().padStart(3, '0')}`;
}

function levelToString(level) {
  switch (level) {
    case LogLevel.DEBUG:
      return 'DEBUG';
    case LogLevel.INFO:
      return 'INFO';
    case LogLevel.WARN:
      return 'WARN';
    case LogLevel.ERROR:
      return 'ERROR';
    default:
      return '';
  }
}

export function write(level, message) {
  const formattedMessage = formatMessage(message);
  consoleLogger.write(level, `${formatDateTime(new Date())}\t[${levelToString(level)} - FRONTEND]\t${formattedMessage}`);
  fileLogger.write(level, `${formatDateTime(new Date())}\t[${levelToString(level)} - FRONTEND]\t${formattedMessage}`);
}

export function setDebug(value) {
  consoleLogger.minLevel = value ? LogLevel.DEBUG : LogLevel.INFO;
}
