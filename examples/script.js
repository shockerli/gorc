#!/usr/bin/env node

const readline = require("readline");
const fs = require('fs')

const rl = readline.createInterface({
  input: process.stdin
});

// listen and read line from stdin
rl.on("line", (input) => {
  fs.writeFile(__dirname + '/data.log', input, () => {
  })
});
