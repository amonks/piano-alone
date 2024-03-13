import * as ts from "typescript";
import * as U from "./util";
import * as fs from "fs";

export function buildAST(names: U.Names): U.AST {
  const programText = `export const _: ${names.typescriptType} = {} as ${names.typescriptType}`;
  fs.writeFileSync("dummy.ts", programText);
  const program = ts.createProgram(["dummy.ts"], {});
  const file = program.getSourceFile("dummy.ts");
  const checker = program.getTypeChecker();
  fs.unlinkSync("dummy.ts");

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
        const name = prop.getName();
        methods[name] = parseMethodDeclaration(name, declarations);
        break;
      default:
        U.debugPrint(
          names,
          "unsupported syntax for",
          prop.getName(),
          ts.SyntaxKind[declarations[0].kind],
        );
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

  function parseMethodDeclaration(
    name: string,
    decs: ts.Declaration[],
  ): U.Arity[] {
    return decs.map((dec) => {
      const children = dec.getChildren();
      if (children[0].kind === ts.SyntaxKind.JSDoc) children.shift();
      const parametersIndex =
        children.findIndex((c) => c.kind === ts.SyntaxKind.OpenParenToken) + 1;
      const returnTypeIndex =
        children.findIndex((c) => c.kind === ts.SyntaxKind.SemicolonToken) - 1;
      const parameters = children[parametersIndex];
      const returnType = children[returnTypeIndex];

      if (parameters.kind !== ts.SyntaxKind.SyntaxList) {
        throw Error("couldn't find parameters");
      }

      return {
        parameters: parseParameters(name, parameters as ts.SyntaxList),
        return: parseReturnType(returnType),
      };
    });
  }

  function parseParameters(name: string, list: ts.SyntaxList): U.Parameter[] {
    const parameters: ts.ParameterDeclaration[] = list
      .getChildren()
      .filter(ts.isParameter);
    if (name === "insertBefore") {
      console.log(
        name,
        parameters.length,
        ts.SyntaxKind[list.getChildren()[0].kind],
      );
    }
    return parameters.map((p) => {
      return {
        name: p.name.getText(),
        type: p.type ? getType(p.type) : U.Literal.Void,
      };
    });
  }

  function getTypeNodeType(node: ts.TypeNode): U.Type {
    const t = checker.getTypeFromTypeNode(node);
    return getType(t);
  }

  function getType(t: ts.Type): U.Type {
    if (t.isUnion()) {
      const name = node.getText();
      const baseTypes = t.getBaseTypes()!;
      // return new U.Union(baseTypes.map((t) => t).map(getType));

      if (name.includes("|")) return checker.typeToString(t);
      if (name === "boolean") return "boolean";
      if (unions[name]) return name;
      if (checker.typeToString(t.types[0])[0] !== `"`) return name;
      unions[name] = t.types.map((t) => checker.typeToString(t));
      return name;
    }
    return checker.typeToString(t);
  }

  function getNodeType(syntax: ts.Node): U.Type {
    switch (syntax.kind) {
      case ts.SyntaxKind.VoidKeyword:
        return U.Type.Void;
      case ts.SyntaxKind.TypeReference:
        const t = checker.getTypeFromTypeNode(syntax as ts.TypeNode);
        return getType(t);
      case ts.SyntaxKind.UnionType:
        return U.Type.fromNode(checker, unions, syntax);
      default:
        return U.Type.fromString(ts.SyntaxKind[syntax.kind]);
    }
  }
}
