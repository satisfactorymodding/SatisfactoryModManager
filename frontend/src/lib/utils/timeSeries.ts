export interface TimeSeries {
  addValue: (value: number) => void;
  getAverage: () => number;
  getDerivative: () => number | undefined;
  clear: () => void;
}

export function timeSeries(millisecondsLifetime: number): TimeSeries {
  const items: { value: number, timestamp: number }[] = [];
  return {
    addValue: (value: number) => {
      items.push({ value, timestamp: Date.now() });
      setTimeout(() => {
        items.shift();
      }, millisecondsLifetime);
    },
    getAverage: () => {
      return items.reduce((a, b) => a + b.value, 0) / items.length;
    },
    getDerivative: () => {
      return (items[items.length - 1].value - items[0].value) / ((items[items.length - 1].timestamp - items[0].timestamp) / 1000); // per second
    },
    clear: () => {
      items.length = 0;
    },
  };
}
