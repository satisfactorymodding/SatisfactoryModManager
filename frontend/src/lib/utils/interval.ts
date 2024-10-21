export function setIntervalImmediate(handler: () => void, timeout: number): ReturnType<typeof setInterval> {
  handler();
  return setInterval(handler, timeout);
}