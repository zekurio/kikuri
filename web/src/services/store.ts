import { AppTheme, getSystemTheme } from "../theme/theme";

import LocalStorageUtil from "../util/localstorage";
import { create } from "zustand";

export type FetchLocked<T> = {
  value: T | undefined;
  isFetching: boolean;
};

export interface Store {
  theme: AppTheme;
  setTheme: (v: AppTheme) => void;

  accentColor?: string;
  setAccentColor: (v?: string) => void;
}

export const useStore = create<Store>((set) => ({
  theme: LocalStorageUtil.get("kikuri.theme", getSystemTheme())!,
  setTheme: (theme) => {
    set({ theme });
    LocalStorageUtil.set("kikuri.theme", theme);
  },

  accentColor: LocalStorageUtil.get("kikuri.accentcolor"),
  setAccentColor: (accentColor) => {
    set({ accentColor });
    if (accentColor === undefined) LocalStorageUtil.del("kikuri.accentcolor");
    else LocalStorageUtil.set("kikuri.accentcolor", accentColor);
  },
}));
