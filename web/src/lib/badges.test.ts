import { describe, expect, it } from 'vitest';
import { badgeClass, componentTone, instrumentTone, semconvTone } from './badges';

describe('instrumentTone', () => {
  it.each([
    ['counter', 'blue'],
    ['updowncounter', 'cyan'],
    ['gauge', 'green'],
    ['histogram', 'orange'],
    ['summary', 'purple'],
    ['mystery', 'neutral'],
  ])('%s is %s', (type, tone) => {
    expect(instrumentTone(type)).toBe(tone);
  });
});

describe('semconvTone', () => {
  it('exact is filled yellow', () => {
    expect(semconvTone('exact')).toEqual({ tone: 'yellow', outlined: false });
  });
  it('prefix is outlined yellow', () => {
    expect(semconvTone('prefix')).toEqual({ tone: 'yellow', outlined: true });
  });
  it('none has no badge', () => {
    expect(semconvTone('none')).toBeNull();
    expect(semconvTone('')).toBeNull();
    expect(semconvTone(undefined)).toBeNull();
  });
});

describe('componentTone', () => {
  it.each(['receiver', 'processor', 'exporter', 'extension', 'connector', 'platform'])('%s is neutral', (type) => {
    expect(componentTone(type)).toBe('neutral');
  });
});

describe('badgeClass', () => {
  it('neutral uses secondary background and weak border', () => {
    expect(badgeClass('neutral')).toContain('bg-secondary');
    expect(badgeClass('neutral')).toContain('border-weak');
  });
  it('a hue tints its background and colours its text', () => {
    expect(badgeClass('blue')).toContain('text-hue-blue');
    expect(badgeClass('blue')).toContain('bg-hue-blue/15');
  });
  it('outlined drops the tint and draws a hue border', () => {
    expect(badgeClass('yellow', true)).toContain('border-hue-yellow');
    expect(badgeClass('yellow', true)).not.toContain('bg-hue-yellow/15');
  });
});
