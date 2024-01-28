import * as ts from "typescript";

function main() {
  const program = ts.createProgram(["dummy.ts"], {});
  const file = program.getSourceFile("dummy.ts");
  const checker = program.getTypeChecker();

  if (!file) throw Error("no file");

  let typeNode: ts.Node | null = null;
  findTypeNode(file);
  if (!typeNode) throw Error("no type node");

  const properties: Record<string, Property> = {};
  const methods: Record<string, Arity[]> = {};
  const unions: Record<string, string[]> = {};

  const t = checker.getTypeAtLocation(typeNode);
  for (const prop of t.getProperties()) {
    const declarations = prop.getDeclarations();
    if (!declarations) throw Error("no declarations");
    switch (declarations[0].kind) {
      case ts.SyntaxKind.PropertySignature:
        properties[prop.getName()] = parsePropertyDeclaration(declarations[0]);
        break;
      case ts.SyntaxKind.MethodSignature:
        methods[prop.getName()] = parseMethodDeclaration(declarations);
        break;
    }
  }

  console.log(`package c2d`);
  console.log(``);
  console.log(`import "syscall/js"`);
  console.log(``);
  console.log(`type C2D js.Value`);
  console.log(``);
  printProperties();
  printMethods();
  printUnions();

  function printProperties() {
    for (const [name, type] of Object.entries(properties)) {
      printProperty(name, type);
    }
  }

  function printProperty(name: string, type: string) {
    console.log(`func (c2d C2D) Get${toPascal(name)}() ${resolve(type)} {`);
    console.log(`	return js.Value(c2d).Get("${name}")${getter(type)}`);
    console.log(`}`);
    console.log(``);
    console.log(`func (c2d C2D) Set${toPascal(name)}(val ${resolve(type)}) {`);
    console.log(`	js.Value(c2d).Set("${name}", val)`);
    console.log(`}`);
    console.log(``);
  }

  function printMethods() {
    for (const [name, arities] of Object.entries(methods)) {
      printArity(name, arities[0]);
    }
  }

  function printArity(name: string, arity: Arity) {
    console.log(
      `func (c2d C2D) ${toPascal(name)}(${plist(arity.parameters)}) ${resolve(
        arity.return,
      )} {`,
    );
    printCall(name, arity);
    console.log(`}`);
    console.log(``);
  }

  function printCall(name: string, arity: Arity) {
    const prefix = arity.return === "void" ? "" : "return ";
    let params = "";
    for (const p of arity.parameters) {
      params += ", " + p.name;
    }
    console.log(`	${prefix}js.Value(c2d).Call("${name}"${params})`);
  }

  function plist(parameters: Parameter[]): string {
    let out: string[] = [];
    for (let i = 0; i < parameters.length; i++) {
      const next = parameters[i + 1];
      const p = parameters[i];
      let s = p.name;
      if (!next || next.type !== p.type) {
        s += " " + resolve(p.type);
      }
      out.push(s);
    }
    return out.join(", ");
  }

  function getter(type: string): string {
    if (type === "string") {
      return `.String()`;
    } else if (type === "boolean") {
      return `.Bool()`;
    } else if (unions[type]) {
      return `.String()`;
    } else if (type === "number") {
      return `.Float()`;
    } else {
      return ``;
    }
  }

  function resolve(name: string): string {
    if (unions[name]) return name;
    if (name === "void") return "";
    if (name === "number") return "float64";
    if (name === "string") return "string";
    if (name === "boolean") return "bool";
    if (name === "boolean | undefined") return "bool";
    return `js.Value`;
  }

  function printUnions() {
    for (const [name, values] of Object.entries(unions)) {
      printUnion(name, values);
    }
  }

  function printUnion(name: string, values: string[]) {
    console.log(`type ${name} = string`);
    console.log(`const (`);
    for (const val of values) {
      console.log(`  ${name}${toPascal(unquote(val))} ${name} = ${val}`);
    }
    console.log(")");
    console.log("");
  }

  function findTypeNode(n: ts.Node) {
    if (typeNode) return;
    if (ts.isTypeNode(n)) typeNode = n;
    ts.forEachChild(n, findTypeNode);
  }

  function parsePropertyDeclaration(dec: ts.Declaration): Property {
    const children = dec.getChildren();
    const propTypeDec = children[1 + children.findIndex(ts.isColonToken)];
    return getType(propTypeDec);
  }

  function parseMethodDeclaration(decs: ts.Declaration[]): Arity[] {
    return decs.map((dec) => {
      const children = dec.getChildren();
      if (children[0].kind === ts.SyntaxKind.JSDoc) children.shift();
      const parameters = parseParameters(children[2] as ts.SyntaxList);
      const returnType = parseReturnType(children[5]);
      return {
        parameters: parameters,
        return: returnType,
      };
    });
  }

  function parseReturnType(syntax: ts.Node): string {
    switch (syntax.kind) {
      case ts.SyntaxKind.VoidKeyword:
        return "void";
      case ts.SyntaxKind.TypeReference:
        const t = checker.getTypeFromTypeNode(syntax as ts.TypeNode);
        return checker.typeToString(t);
      case ts.SyntaxKind.UnionType:
        return getType(syntax);
      default:
        return ts.SyntaxKind[syntax.kind];
    }
  }

  function parseParameters(list: ts.SyntaxList): Parameter[] {
    const parameters: ts.ParameterDeclaration[] = list
      .getChildren()
      .filter(ts.isParameter);
    return parameters.map((p) => {
      return {
        name: p.name.getText(),
        type: p.type ? getType(p.type) : "void",
      };
    });
  }

  function getType(node: ts.Node): string {
    const t = checker.getTypeFromTypeNode(node as ts.TypeNode);
    if (t.isUnion()) {
      const name = node.getText();
      if (name.includes("|")) return checker.typeToString(t);
      if (name === "boolean") return "boolean";
      if (unions[name]) return name;
      if (checker.typeToString(t.types[0])[0] !== `"`) return name;
      unions[name] = t.types.map((t) => checker.typeToString(t));
      return name;
    }
    return checker.typeToString(t);
  }
}

function unquote(text: string) {
  return String(JSON.parse(text));
}

function toPascal(text: string) {
  return text.replace(/(^\w|-\w)/g, clearAndUpper);
}

function clearAndUpper(text: string) {
  return text.replace(/-/, "").toUpperCase();
}

type Property = string;
type Arity = {
  parameters: Parameter[];
  return: string;
};
type Parameter = { name: string; type: string };

main();
