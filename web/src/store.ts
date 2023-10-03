import { create } from 'zustand'

type Store = {

  loggedIn: boolean;
  setLoggedIn: (loggedIn: boolean) => void;

  isAdmin: boolean;
  setIsAdmin: (isAdmin: boolean) => void;
};

export const useStore = create<Store>((set, get) => ({
  loggedIn: false,
  setLoggedIn: (loggedIn) => set({ loggedIn }),

  isAdmin: false,
  setIsAdmin: (isAdmin) => set({ isAdmin }),
}));