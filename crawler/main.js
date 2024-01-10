const { argv } = require("node:process");
const { crawlPage } = require("./crawl");
const { printReport } = require("./report");

function main() {
  if (argv.length < 3) {
    console.error("You need to pass an argument");
    return;
  }
  if (argv.length > 3) {
    console.error("You passed too many arguments");
    return;
  }
  const url = argv[2];
  crawlPage(url, url, {}).then(pages => {
    printReport(pages);
  });
}

main();