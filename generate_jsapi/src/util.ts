import * as fs from "fs";
import * as ts from "typescript";

export type AST = {
  properties: Record<string, Property>;
  methods: Record<string, Arity[]>;
  unions: Record<string, string[]>;
};
export type Property = string;
export type Arity = {
  parameters: Parameter[];
  return: Type;
  suffix?: string;
};
export type Parameter = {
  name: string;
  type: Type;
};
export type Names = {
  typescriptType: string;
  goLocal: string;
  goType: string;
  goFile: string;
  goPackage: string;
};

export interface Type {
  toString(): string;
  explode(): Type[];
}

export class Literal implements Type {
  constructor(readonly name: string) {}
  toString() {
    return this.name;
  }
  explode() {
    return [this];
  }
  static Void = new Literal("void");
  static Undefined = new Literal("undefined");
  static Null = new Literal("null");
  static Boolean = new Literal("boolean");
  static String = new Literal("string");
  static Float = new Literal("float");
}

export class StringLiteral implements Type {
  constructor(readonly s: string) {}
  toString() {
    return `"${this.s}"`;
  }
  explode() {
    return [this];
  }
}

export class Reference implements Type {
  constructor(readonly name: string) {}
  toString() {
    return this.name;
  }
  explode() {
    return [this];
  }
}

export class Union implements Type {
  constructor(readonly values: Type[]) {}
  toString() {
    return this.values.map((v) => v.toString).join(" | ");
  }
  explode() {
    return this.values.reduce(
      (acc, v) => [...acc, ...v.explode()],
      [] as Type[],
    );
  }
}

export class Optional implements Type {
  constructor(readonly value: Type) {}
  toString() {
    return this.value.toString() + "?";
  }
  explode() {
    return [Literal.Undefined, this.value];
  }
}

export function unquote(text: string) {
  return String(JSON.parse(text));
}

export function toPascal(text: string) {
  return text.replace(/(^\w|-\w)/g, clearAndUpper);
}

export function clearAndUpper(text: string) {
  return text.replace(/-/, "").toUpperCase();
}

export function debugPrint(names: Names, ...values: any[]) {
  if (process.env.DEBUG_GENERATE_JSAPI !== names.goFile) {
    return;
  }
  console.log(...values);
}

export function debugWrite(name: string, names: Names, ast: AST) {
  if (process.env.DEBUG_GENERATE_JSAPI !== names.goFile) {
    return;
  }
  fs.writeFileSync(name + ".json", JSON.stringify(ast, null, 2));
}
