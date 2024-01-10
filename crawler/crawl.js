const { JSDOM } = require('jsdom');

function normalizeURL(url) {
  const urlObj = new URL(url);
  urlObj.hash = '';
  urlObj.search = '';
  urlObj.protocol = '';
  urlObj.host = '';
  let res = urlObj.toString();
  if (res.endsWith('/')) res = res.slice(0, -1);
  res = res.replace("https://", '');
  res = res.replace("http://", '')
  return res;
}

function getURLsFromHTML(html, baseURL) {
  const dom = new JSDOM(html);
  const unparsedURLs = [...dom.window.document.querySelectorAll("a")].map(a => a.href);
  const urls = unparsedURLs.map(url => url.startsWith("http") ? url : new URL(url, baseURL).toString());
  return urls;
}

function isSameDomain(url1, url2) {
  const url1Obj = new URL(url1);
  const url2Obj = new URL(url2);
  return url1Obj.hostname === url2Obj.hostname;
}

async function crawlPage(baseURL, currentURL, pages) {
  if (!isSameDomain(baseURL, currentURL)) {
    return pages;
  }
  const normalized = normalizeURL(currentURL);
  console.log(`Crawling ${normalized}...`);
  if (normalized in pages) {
    pages[normalized] += 1;
    return pages;
  }
  pages[normalized] = baseURL !== currentURL ? 1 : 0;
  try {
    const response = await fetch(currentURL);
    if (response.status >= 400) {
      return pages;
    }
    if (!response.headers.get("Content-Type").includes("text/html")) {
      return pages;
    }
    const html = await response.text();
    const urls = getURLsFromHTML(html, baseURL);
    for (let i = 0; i < urls.length; i++) {
      const url = urls[i];
      pages = await crawlPage(baseURL, url, pages);
    }
  } catch (error) {
    console.error(`Error fetching ${baseURL}: ${error}`);
    return;
  }
  console.log(pages)
  return pages;
}

module.exports = {
  normalizeURL,
  getURLsFromHTML,
  isSameDomain,
  crawlPage
}