export type AST = {
  properties: Record<string, Property>;
  methods: Record<string, Arity[]>;
  unions: Record<string, string[]>;
};
export type Property = string;
export type Arity = {
  parameters: Parameter[];
  return: string;
  suffix?: string;
};
export type Parameter = { name: string; type: string };

export function unquote(text: string) {
  return String(JSON.parse(text));
}

export function toPascal(text: string) {
  return text.replace(/(^\w|-\w)/g, clearAndUpper);
}

export function clearAndUpper(text: string) {
  return text.replace(/-/, "").toUpperCase();
}
