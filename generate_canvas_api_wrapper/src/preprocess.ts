import * as U from "./util";

export function explodeAST(ast: U.AST): U.AST {
  const ret: U.AST = { ...ast, methods: {} };
  for (const [name, arities] of Object.entries(ast.methods)) {
    ret.methods[name] = explodeArities(arities);
  }
  return ret;
}

export function dedupeAST(ast: U.AST): U.AST {
  const ret: U.AST = { ...ast, methods: {} };
  for (const [name, arities] of Object.entries(ast.methods)) {
    ret.methods[name] = dedupeArities(arities);
  }
  return ret;
}

function explodeArities(arities: U.Arity[]): U.Arity[] {
  const out: U.Arity[] = [];
  for (const arity of arities) {
    const firstOptionalIndex = arity.parameters.findLastIndex((p) =>
      p.type.endsWith("?"),
    );
    if (firstOptionalIndex === -1) {
      out.push(arity);
      continue;
    }
    for (let i = firstOptionalIndex; i < arity.parameters.length; i++) {
      out.push({
        ...arity,
        parameters: arity.parameters
          .slice(0, i)
          .map((p) => ({ name: p.name, type: removeSuffix(p.type, "?") })),
      });
    }
    out.push({
      ...arity,
      parameters: arity.parameters
        .slice(0)
        .map((p) => ({ name: p.name, type: removeSuffix(p.type, "?") })),
    });
  }
  return dedupeArities(out);
}

function dedupeArities(arities: U.Arity[]): U.Arity[] {
  arities.sort((a, b) => {
    const ldiff = a.parameters.length - b.parameters.length;
    if (ldiff !== 0) return ldiff;
    return JSON.stringify(a.parameters) < JSON.stringify(b.parameters) ? -1 : 1;
  });
  const out: U.Arity[] = [];
  let [prev, ...rest] = arities;
  out.push(prev);
  for (const arity of rest) {
    if (identifyArity(arity) === identifyArity(prev)) {
      continue;
    }
    if (arity.parameters.length === prev.parameters.length) {
      out.pop();
      out.push({
        ...prev,
        suffix: U.toPascal(prev.parameters.map((p) => p.type).join("-")),
      });
      out.push({
        ...arity,
        suffix: U.toPascal(arity.parameters.map((p) => p.type).join("-")),
      });
      prev = arity;
    } else {
      prev = arity;
      out.push(arity);
    }
  }
  return out;
}

function identifyArity(arity: U.Arity): string {
  return arity.parameters.map((p) => p.type).join(",");
}

function removeSuffix(s: string, suffix: string): string {
  if (s.endsWith(suffix)) {
    return s.slice(0, s.length - suffix.length);
  }
  return s;
}
