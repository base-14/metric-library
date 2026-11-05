export type BadgeTone = 'blue' | 'green' | 'orange' | 'purple' | 'cyan' | 'yellow' | 'neutral';

const INSTRUMENT_TONES: Record<string, BadgeTone> = {
  counter: 'blue',
  updowncounter: 'cyan',
  gauge: 'green',
  histogram: 'orange',
  summary: 'purple',
};

const FILLED: Record<BadgeTone, string> = {
  blue: 'bg-hue-blue/15 text-hue-blue border-transparent',
  green: 'bg-hue-green/15 text-hue-green border-transparent',
  orange: 'bg-hue-orange/15 text-hue-orange border-transparent',
  purple: 'bg-hue-purple/15 text-hue-purple border-transparent',
  cyan: 'bg-hue-cyan/15 text-hue-cyan border-transparent',
  yellow: 'bg-hue-yellow/15 text-hue-yellow border-transparent',
  neutral: 'bg-secondary text-secondary border-weak',
};

const OUTLINED: Record<BadgeTone, string> = {
  blue: 'bg-transparent text-hue-blue border-hue-blue',
  green: 'bg-transparent text-hue-green border-hue-green',
  orange: 'bg-transparent text-hue-orange border-hue-orange',
  purple: 'bg-transparent text-hue-purple border-hue-purple',
  cyan: 'bg-transparent text-hue-cyan border-hue-cyan',
  yellow: 'bg-transparent text-hue-yellow border-hue-yellow',
  neutral: 'bg-transparent text-secondary border-weak',
};

export function instrumentTone(type: string): BadgeTone {
  return INSTRUMENT_TONES[type] ?? 'neutral';
}

export function semconvTone(match?: string): { tone: BadgeTone; outlined: boolean } | null {
  if (match === 'exact') return { tone: 'yellow', outlined: false };
  if (match === 'prefix') return { tone: 'yellow', outlined: true };
  return null;
}

const COMPONENT_TONES: Record<string, BadgeTone> = {};

export function componentTone(type: string): BadgeTone {
  return COMPONENT_TONES[type] ?? 'neutral';
}

export function badgeClass(tone: BadgeTone, outlined = false): string {
  const colours = outlined ? OUTLINED[tone] : FILLED[tone];
  return `inline-flex items-center px-1.5 py-0.5 text-xs font-medium rounded-[4px] border ${colours}`;
}
