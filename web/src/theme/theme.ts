import Color from "color";

export enum AppTheme {
  DARK = 0,
  LIGHT = 1,
}

export const DarkTheme = {
  background: "#11111b",
  background2: "#181825",
  background3: "#1e1e2e",

  text: "#cdd6f4",
  textAlt: "#4c4f69",

  accent: "#88b3f9",
  accentDarker: Color.xyz("#618dbd").darken(0.3).hexa(),

  rosewater: "#f5e0dc",
  flamingo: "#f2cdcd",
  pink: "#f5c2e7",
  mauve: "#cba6f7",
  red: "#f38ba8",
  maroon: "#eba0ac",
  peach: "#fab387",
  yellow: "#f9e2af",
  green: "#a6e3a1",
  teal: "#94e2d5",
  sky: "#89dceb",
  sapphire: "#74c7ec",
  blue: "#89b4fa",
  lavender: "#b4befe",

  _isDark: true,
};

export const LightTheme: Theme = {
  ...DarkTheme,

  background3: "#eff1f5",
  background2: "#e6e9ef",
  background: "#dddddd",

  text: "#4c4f69",
  textAlt: "#cdd6f4",

  rosewater: "#dc8a78",
  flamingo: "#dd7878",
  pink: "#ea76cb",
  mauve: "#8839ef",
  red: "#d20f39",
  maroon: "#e64553",
  peach: "#fe640b",
  yellow: "#df8e1d",
  green: "#40a02b",
  teal: "#179299",
  sky: "#04a5e5",
  sapphire: "#209fb5",
  blue: "#1e66f5",
  lavender: "#7287fd",

  _isDark: false,
};

export const DefaultTheme = DarkTheme;
export type Theme = typeof DefaultTheme;

export const getSystemTheme = () => {
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? AppTheme.DARK
    : AppTheme.LIGHT;
};
