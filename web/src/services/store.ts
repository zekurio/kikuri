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
  theme: LocalStorageUtil.get("daemon.theme", getSystemTheme())!,
  setTheme: (theme) => {
    set({ theme });
    LocalStorageUtil.set("daemon.theme", theme);
  },

  accentColor: LocalStorageUtil.get("daemon.accentcolor"),
  setAccentColor: (accentColor) => {
    set({ accentColor });
    if (accentColor === undefined) LocalStorageUtil.del("daemon.accentcolor");
    else LocalStorageUtil.set("daemon.accentcolor", accentColor);
  },
}));
