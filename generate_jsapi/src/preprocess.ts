import * as U from "./util";

export function explodeAST(ast: U.AST): U.AST {
  const ret: U.AST = { ...ast, methods: {} };
  for (const [name, arities] of Object.entries(ast.methods)) {
    let exploded = arities;
    exploded = explodeOptionals(exploded);
    exploded = explodeUnions(exploded);
    exploded = explodeOptionals(exploded);
    ret.methods[name] = exploded;
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

function explodeOptionals(arities: U.Arity[]): U.Arity[] {
  const out: U.Arity[] = [];
  for (const arity of arities) {
    const firstOptionalIndex = arity.parameters.findLastIndex((p) =>
      p.type.isOptional(),
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
          .map((p) => ({ name: p.name, type: p.type.toRequired() })),
      });
    }
    out.push({
      ...arity,
      parameters: arity.parameters
        .slice(0)
        .map((p) => ({ name: p.name, type: p.type.toRequired() })),
    });
  }
  return dedupeArities(out);
}

function explodeUnions(arities: U.Arity[]): U.Arity[] {
  const out: U.Arity[] = [];
  for (const arity of arities) {
    recur(arity.parameters);

    function recur(parameters: U.Parameter[]) {
      if (parameters.every((p) => !p.type.includes("|"))) {
        out.push({ ...arity, parameters });
        return;
      }

      const index = parameters.findIndex((p) => p.type.includes("|"));
      const values = splitUnion(parameters[index].type);
      for (const value of values) {
        recur([
          ...parameters.slice(0, index),
          { ...parameters[index], type: value },
          ...parameters.slice(index + 1),
        ]);
      }
    }
  }
  return dedupeArities(out);
}

function splitUnion(union: string): string[] {
  const out: string[] = [];
  let partial = "";
  let depth = 0;
  let recur = false;
  for (const char of union) {
    switch (char) {
      case "|":
        if (depth === 0) {
          out.push(partial);
          partial = "";
        } else {
          partial += "|";
          recur = true;
        }
        break;

      case "(":
        if (depth > 0) {
          partial += ")";
        }
        depth += 1;
        break;

      case ")":
        depth -= 1;
        if (depth > 0) {
          partial += ")";
        }
        out.push(partial);
        partial = "";
        break;

      case " ":
        break;

      default:
        partial += char;
        break;
    }
  }
  if (recur) {
    return out.reduce((acc, u) => [...acc, ...splitUnion(u)], [] as string[]);
  }
  return out;
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
