import * as U from "./util";

class StringBuilder {
  str = "";
  add(s: string) {
    this.str += s + "\n";
  }
  string(): string {
    return this.str;
  }
}

export function printAST(ast: U.AST): string {
  const s = new StringBuilder();
  s.add(printHeader());
  s.add(printProperties(ast));
  s.add(printMethods(ast));
  s.add(printUnions(ast));
  return s.string();
}

function printHeader(): string {
  const s = new StringBuilder();
  s.add(`package c2d`);
  s.add(``);
  s.add(`import "syscall/js"`);
  s.add(``);
  s.add(`type C2D js.Value`);
  s.add(``);
  return s.string();
}

function printUnions(ast: U.AST): string {
  const s = new StringBuilder();
  for (const [name, values] of Object.entries(ast.unions)) {
    s.add(printUnion(name, values));
  }
  return s.string();
}

function printUnion(name: string, values: string[]): string {
  const s = new StringBuilder();
  s.add(`type ${name} = string`);
  s.add(`const (`);
  for (const val of values) {
    s.add(`  ${name}${U.toPascal(U.unquote(val))} ${name} = ${val}`);
  }
  s.add(")");
  s.add("");
  return s.string();
}

function printProperties(ast: U.AST): string {
  const s = new StringBuilder();
  for (const [name, type] of Object.entries(ast.properties)) {
    s.add(printProperty(ast, name, type));
  }
  return s.string();
}

function printProperty(ast: U.AST, name: string, type: string): string {
  const s = new StringBuilder();
  s.add(
    `func (c2d C2D) Get${U.toPascal(name)}() ${resolveOutput(ast, type)} {`,
  );
  s.add(`	return js.Value(c2d).Get("${name}")${printGetter(ast, type)}`);
  s.add(`}`);
  s.add(``);
  s.add(
    `func (c2d C2D) Set${U.toPascal(name)}(val ${resolveInput(ast, type)}) {`,
  );
  s.add(`	js.Value(c2d).Set("${name}", val)`);
  s.add(`}`);
  s.add(``);
  return s.string();
}

function printMethods(ast: U.AST): string {
  const s = new StringBuilder();
  for (const [name, arities] of Object.entries(ast.methods)) {
    for (const arity of arities) {
      s.add(printArity(ast, name, arity, arities.length > 1));
    }
  }
  return s.string();
}

function printArity(
  ast: U.AST,
  name: string,
  arity: U.Arity,
  useOrdinal: boolean,
): string {
  const s = new StringBuilder();

  const ordinal = useOrdinal ? String(arity.parameters.length) : "";
  const fullName = U.toPascal(name) + ordinal + (arity.suffix || "");
  const params = printParamList(ast, arity.parameters);
  const ret = resolveOutput(ast, arity.return);
  s.add(`func (c2d C2D) ${fullName}(${params}) ${ret} {`);
  s.add(printCall(name, arity));
  s.add(`}`);
  s.add(``);
  return s.string();
}

function printCall(name: string, arity: U.Arity): string {
  const prefix = arity.return === "void" ? "" : "return ";
  let params = "";
  for (const p of arity.parameters) {
    params += ", " + p.name;
  }
  return `	${prefix}js.Value(c2d).Call("${name}"${params})`;
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
