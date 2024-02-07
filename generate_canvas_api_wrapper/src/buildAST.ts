import * as ts from "typescript";
import * as U from "./util";

export function buildAST(): U.AST {
  const program = ts.createProgram(["dummy.ts"], {});
  const file = program.getSourceFile("dummy.ts");
  const checker = program.getTypeChecker();

  if (!file) throw Error("no file");

  let typeNode: ts.Node | null = null;
  findTypeNode(file);
  if (!typeNode) throw Error("no type node");

  const properties: Record<string, U.Property> = {};
  const methods: Record<string, U.Arity[]> = {};
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

  return { properties, methods, unions };

  function findTypeNode(n: ts.Node) {
    if (typeNode) return;
    if (ts.isTypeNode(n)) typeNode = n;
    ts.forEachChild(n, findTypeNode);
  }

  function parsePropertyDeclaration(dec: ts.Declaration): U.Property {
    const children = dec.getChildren();
    const propTypeDec = children[1 + children.findIndex(ts.isColonToken)];
    return getType(propTypeDec);
  }

  function parseMethodDeclaration(decs: ts.Declaration[]): U.Arity[] {
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

  function parseParameters(list: ts.SyntaxList): U.Parameter[] {
    const parameters: ts.ParameterDeclaration[] = list
      .getChildren()
      .filter(ts.isParameter);
    return parameters.map((p) => {
      const parsed = {
        name: p.name.getText(),
        type: p.type ? getType(p.type) : "void",
      };
      if (p.questionToken) {
        parsed.type += "?";
      }
      return parsed;
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
}
