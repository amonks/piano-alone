import * as U from "./util";

export function printAST(ast: U.AST) {
  console.log(`package c2d`);
  console.log(``);
  console.log(`import "syscall/js"`);
  console.log(``);
  console.log(`type C2D js.Value`);
  console.log(``);
  printProperties(ast);
  printMethods(ast);
  printUnions(ast);
}

function printUnions(ast: U.AST) {
  for (const [name, values] of Object.entries(ast.unions)) {
    printUnion(name, values);
  }
}

function printUnion(name: string, values: string[]) {
  console.log(`type ${name} = string`);
  console.log(`const (`);
  for (const val of values) {
    console.log(`  ${name}${U.toPascal(U.unquote(val))} ${name} = ${val}`);
  }
  console.log(")");
  console.log("");
}

function printProperties(ast: U.AST) {
  for (const [name, type] of Object.entries(ast.properties)) {
    printProperty(ast, name, type);
  }
}

function printProperty(ast: U.AST, name: string, type: string) {
  console.log(
    `func (c2d C2D) Get${U.toPascal(name)}() ${resolveOutput(ast, type)} {`,
  );
  console.log(`	return js.Value(c2d).Get("${name}")${printGetter(ast, type)}`);
  console.log(`}`);
  console.log(``);
  console.log(
    `func (c2d C2D) Set${U.toPascal(name)}(val ${resolveInput(ast, type)}) {`,
  );
  console.log(`	js.Value(c2d).Set("${name}", val)`);
  console.log(`}`);
  console.log(``);
}

function printMethods(ast: U.AST) {
  for (const [name, arities] of Object.entries(ast.methods)) {
    for (const arity of arities) {
      printArity(ast, name, arity, arities.length > 1);
    }
  }
}

function printArity(
  ast: U.AST,
  name: string,
  arity: U.Arity,
  useOrdinal: boolean,
) {
  const ordinal = useOrdinal ? String(arity.parameters.length) : "";
  const fullName = U.toPascal(name) + ordinal + (arity.suffix || "");
  const params = printParamList(ast, arity.parameters);
  const ret = resolveOutput(ast, arity.return);
  console.log(`func (c2d C2D) ${fullName}(${params}) ${ret} {`);
  printCall(name, arity);
  console.log(`}`);
  console.log(``);
}

function printCall(name: string, arity: U.Arity) {
  const prefix = arity.return === "void" ? "" : "return ";
  let params = "";
  for (const p of arity.parameters) {
    params += ", " + p.name;
  }
  console.log(`	${prefix}js.Value(c2d).Call("${name}"${params})`);
}

function printParamList(ast: U.AST, parameters: U.Parameter[]): string {
  let out: string[] = [];
  for (let i = 0; i < parameters.length; i++) {
    const next = parameters[i + 1];
    const p = parameters[i];
    let s = p.name;
    if (!next || next.type !== p.type) {
      s += " " + resolveInput(ast, p.type);
    }
    out.push(s);
  }
  return out.join(", ");
}

function resolveInput(ast: U.AST, name: string): string {
  const resolved = resolveOutput(ast, name);
  if (resolved === "js.Value") return "any";
  return resolved;
}

function resolveOutput(ast: U.AST, name: string): string {
  if (ast.unions[name]) return name;
  if (name === "void") return "";
  if (name === "number") return "float64";
  if (name === "string") return "string";
  if (name === "boolean") return "bool";
  if (name === "boolean | undefined") return "bool";
  return `js.Value`;
}

function printGetter(ast: U.AST, type: string): string {
  if (type === "string") {
    return `.String()`;
  } else if (type === "boolean") {
    return `.Bool()`;
  } else if (ast.unions[type]) {
    return `.String()`;
  } else if (type === "number") {
    return `.Float()`;
  } else {
    return ``;
  }
}
