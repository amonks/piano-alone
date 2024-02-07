import { buildAST } from "./buildAST";
import { dedupeAST, explodeAST } from "./preprocess";
import { printAST } from "./printAST";

const { writeFileSync } = require("fs");

function main() {
  let ast = buildAST();
  ast = explodeAST(ast);
  ast = dedupeAST(ast);

  const out = printAST(ast);
  writeFileSync("out.go", out);
}

main();
