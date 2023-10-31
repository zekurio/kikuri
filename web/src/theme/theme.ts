import Color from "color";

export enum AppTheme {
  DARK = 0,
  LIGHT = 1,
}

export const DarkTheme = {
  background: "#1e1e2e",
  background2: "#313244",
  background3: "#45475a",

  text: "#cdd6f4",
  textAlt: "#bac2de",

  accent: "#9447dc",
  accentDarker: Color.xyz("#6812bb").darken(0.3).hexa(),

  white: "#f4f4f5",
  whiteDarker: "#dddddd",
  blurple: "#5865f2",
  blurpleDarker: "#4450d6",
  darkGray: "#1e1e1e",
  red: "#ed4245",
  orange: "#f57c00",
  yellow: "#fbc02d",
  green: "#43a047",
  lime: "#57f287",
  cyan: "#03a9f4",
  pink: "#eb459e",

  _isDark: true,
};

export const LightTheme: Theme = {
  ...DarkTheme,

  background3: "rgb(250, 250, 250)",
  background2: "rgb(235, 235, 235)",
  background: "rgb(225, 225, 225)",

  text: "#212121",
  textAlt: "#f4f4f5",

  _isDark: false,
};

export const DefaultTheme = DarkTheme;
export type Theme = typeof DefaultTheme;

export const getSystemTheme = () => {
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? AppTheme.DARK
    : AppTheme.LIGHT;
};
