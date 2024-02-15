export const largeNumberFormat = Intl.NumberFormat(undefined, { notation: 'compact' }).format;

// We cannot use the units mode of NumberFormat, since it is not aware of different names for larger units
// For 1 TB, it uses 1 BB (1 billion bytes), and for 1000 seconds it uses 1Ks (1 thousand seconds)

export function roundWithDecimals(number: number, decimals = 0): number {
  return Math.round(number * (10 ** decimals)) / (10 ** decimals);
}

const sizeRanges = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];

export function bytesToAppropriate(bytes: number): string {
  let rangeNum = 0;
  while (bytes >= 1024 ** (rangeNum + 1)) {
    rangeNum += 1;
  }
  return `${roundWithDecimals(bytes / (1024 ** rangeNum), 2).toFixed(2)} ${sizeRanges[rangeNum]}`;
}

const timeRanges = {
  sec: 1,
  min: 60,
  h: 60 * 60,
  days: 60 * 60 * 24,
};

export function secondsToAppropriate(seconds: number): string {
  const ranges = Object.keys(timeRanges) as (keyof (typeof timeRanges))[];
  let rangeNum = 0;
  while (rangeNum < ranges.length - 1 && seconds >= timeRanges[ranges[rangeNum + 1]]) {
    rangeNum += 1;
  }
  return `${roundWithDecimals(seconds / timeRanges[ranges[rangeNum]], 0)}${ranges[rangeNum]}`;
}
