import { buildAST } from "./buildAST";
import { preprocessAST } from "./preprocess";
import { printAST } from "./printAST";

const { readFileSync, writeFileSync } = require("fs");

function main() {
  let ast = buildAST();
  writeFileSync("ast.json", JSON.stringify(ast, null, 2));

  ast = preprocessAST(ast);
  writeFileSync("processed.json", JSON.stringify(ast, null, 2));

  printAST(ast);
}

main();
