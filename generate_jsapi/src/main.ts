import { buildAST } from "./buildAST";
import * as U from "./util";
import { dedupeAST, explodeAST } from "./preprocess";
import { printAST } from "./printAST";
import * as fs from "fs";

function main() {
  const names: U.Names = {
    typescriptType: process.argv[process.argv.length - 5],
    goLocal: process.argv[process.argv.length - 4],
    goType: process.argv[process.argv.length - 3],
    goFile: process.argv[process.argv.length - 2],
    goPackage: process.argv[process.argv.length - 1],
  };
  console.log(names);

  let ast = buildAST(names);
  U.debugWrite("1", names, ast);
  ast = explodeAST(ast);
  U.debugWrite("2", names, ast);
  ast = dedupeAST(ast);

  const out = printAST(names, ast);
  fs.writeFileSync(`${names.goFile}.go`, out);
}

main();
