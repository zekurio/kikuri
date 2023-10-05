export enum AppTheme {
  DARK = 0,
  LIGHT = 1,
}

export const DarkTheme = {
  background: '#0d1117',
  background2: '#161b22',
  background3: '#21262d',

  text: '#c9d1d9',

  accent: '#58a6ff',

  white: '#c9d1d9',
  whiteDarker: '#b1bac4',
  blurple: '#6f42c1',
  blurpleDarker: '#553982',
  gray: '#8b949e',
  darkGray: '#30363d',
  red: '#d73a49',
  orange: '#e36209',
  yellow: '#f1e05a',
  green: '#28a745',
  lime: '#a5d6a7',
  cyan: '#39c5cf',
  pink: '#d871bd',
};

export const LightTheme = {
  background: '#ffffff',
  background2: '#f6f8fa',
  background3: '#e1e4e8',

  text: '#24292e',

  accent: '#0366d6',

  white: '#24292e',
  whiteDarker: '#586069',
  blurple: '#6f42c1',
  blurpleDarker: '#553982',
  gray: '#6a737d',
  darkGray: '#d1d5da',
  red: '#cb2431',
  orange: '#e36209',
  yellow: '#b08800',
  green: '#28a745',
  lime: '#a5d6a7',
  cyan: '#0366d6',
  pink: '#d871bd',
};

export const DefaultTheme = DarkTheme;
export type Theme = typeof DefaultTheme;
