import { AppTheme, getSystemTheme } from "../theme/theme";
import { Guild, User } from "../lib/kikuri-ts/src";
import { ModalState } from '../hooks/useModal';
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

  selfUser: FetchLocked<User>;
  setSelfUser: (selfUser: Partial<FetchLocked<User>>) => void;

  guilds?: Guild[];
  setGuilds: (guilds?: Guild[]) => void;

  selectedGuild?: Guild;
  setSelectedGuild: (selectedGuild: Guild) => void;

  modal: ModalState<any>;
  setModal: (modal: ModalState<any>) => void;
}

export const useStore = create<Store>((set, get) => ({
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

  selfUser: { value: undefined, isFetching: false },
  setSelfUser: (selfUser: Partial<FetchLocked<User>>) =>
    set({ selfUser: { ...get().selfUser, ...selfUser } }),

  guilds: undefined,
  setGuilds: (guilds) => set({ guilds }),

  selectedGuild: undefined,
  setSelectedGuild: (selectedGuild) => {
    set({ selectedGuild });
    if (!!selectedGuild)
      LocalStorageUtil.set("kikuri.selectedguild", selectedGuild.id);
  },

  modal: { isOpen: false },
  setModal: (modal) => set({ modal }),
}));
