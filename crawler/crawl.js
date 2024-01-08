const { JSDOM } = require('jsdom');

function normalizeURL(url) {
  const urlObj = new URL(url);
  urlObj.hash = '';
  urlObj.search = '';
  urlObj.protocol = '';
  urlObj.host = '';
  let res = urlObj.toString();
  if (res.endsWith('/')) res = res.slice(0, -1);
  res = res.replace("https://", '')
  res = res.replace("http://", '')
  return res;
}

function getURLsFromHTML(html, baseURL) {
  const dom = new JSDOM(html);
  const unparsedURLs = [...dom.window.document.querySelectorAll("a")].map(a => a.href);
  const urls = unparsedURLs.map(url => url.startsWith("http") ? url : new URL(url, baseURL).toString());
  return urls;
}

async function crawlPage(url) {
  console.log(`Crawling ${url}...`);
  try {
    const response = await fetch(url);
    if (response.status >= 400) {
      console.error(`Received status code ${response.status} for ${url}`);
      return;
    }
    if (!response.headers.get("Content-Type").includes("text/html")) {
      console.error(`Expected HTML for ${url} but received ${response.headers.get("Content-Type")}`);
      return;
    }
    const html = await response.text();
    console.log(html);
  } catch (error) {
    console.error(`Error fetching ${url}: ${error}`);
    return;
  }

}

module.exports = {
  normalizeURL,
  getURLsFromHTML,
  crawlPage
}