const { test, expect } = require('@jest/globals');
const { normalizeURL, getURLsFromHTML } = require('./crawl');

describe('normalizeURL', () => {
  test('removes hash', () => {
    const url = 'https://example.com#hash';
    const normalizedURL = normalizeURL(url);
    expect(normalizedURL).not.toContain('#');
  });

  test('removes search', () => {
    const url = 'https://example.com?search=test';
    const normalizedURL = normalizeURL(url);
    expect(normalizedURL).not.toContain('?');
  });

  test('removes hash and search', () => {
    const url = 'https://example.com?search=test#hash';
    const normalizedURL = normalizeURL(url);
    expect(normalizedURL).not.toContain('?');
    expect(normalizedURL).not.toContain('#');
  });

  test('removes protocol', () => {
    const url = 'https://example.com';
    const normalizedURL = normalizeURL(url);
    expect(normalizedURL).not.toContain('https://');
  });

  test('removes trailing slash', () => {
    const url = 'https://example.com/path/';
    const normalizedURL = normalizeURL(url);
    expect(normalizedURL).not.toContain('/path/');
    expect(normalizedURL).toContain('/path');
  });
});

describe('getURLsFromHTML', () => {
  const testHTML = `
    <html>
      <head>
        <link rel="stylesheet" href="https://example.com/style.css">
      </head>
      <body>
        <a href="https://example.com">Home</a>
        <a href="/help">Help</a>
        <a href="https://example.com/about">About</a>
        <a href="https://example.com/contact">Contact</a>
      </body>
    </html>
  `;

  test('returns an array', () => {
    const urls = getURLsFromHTML(testHTML, 'https://example.com');
    expect(Array.isArray(urls)).toBe(true);
  });

  test("returns 4 URLs", () => {
    const urls = getURLsFromHTML(testHTML, 'https://example.com');
    expect(urls.length).toBe(4);
  });

  test("returns absolute URLs if href is relative", () => {
    const urls = getURLsFromHTML(testHTML, 'https://example.com');
    expect(urls[1]).toBe('https://example.com/help');
  });
});