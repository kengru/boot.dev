function printReport(pages) {
  console.log("Starting report...")
  console.log(`Crawled ${Object.keys(pages).length} pages.`);
  const arr = Object.entries(pages).sort(([, a], [, b]) => b - a).reduce((r, [k, v]) => ({ ...r, [k]: v }), {});
  Object.keys(arr).forEach(page => {
    console.log(`Found ${pages[page]} links to: ${page}`);
  });
}

module.exports = {
  printReport
}