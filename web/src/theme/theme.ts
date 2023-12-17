import Color from "color";

export enum AppTheme {
  DARK = 0,
  LIGHT = 1,
}

export const DarkTheme = {
  background: '#111417',
  background2: '#191c20',
  background3: '#25282c',

  text: '#f4f4f5',
  textAlt: '#f4f4f5',

  accent: '#c24de2',
  accentDarker: '#862db3',

  white: '#f4f4f5',
  whiteDarker: '#dddddd',
  blurple: '#5865f2',
  blurpleDarker: '#4450d6',
  darkGray: '#1e1e1e',
  red: '#ed4245',
  orange: '#ea7d0f',
  yellow: '#ffd817',
  green: '#338838',
  lime: '#3ed56f',
  cyan: '#03bcf4',
  pink: '#eb459e',

  _isDark: true,
};

export const LightTheme: Theme = {
  ...DarkTheme,

  background3: 'rgb(250, 250, 250)',
  background2: 'rgb(235, 235, 235)',
  background: 'rgb(225, 225, 225)',

  text: '#212121',
  textAlt: '#f4f4f5',

  accentDarker: '#bd6ffd',

  _isDark: false,
};

export const DefaultTheme = DarkTheme;
export type Theme = typeof DefaultTheme;

export const getSystemTheme = () => {
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? AppTheme.DARK
    : AppTheme.LIGHT;
};