import create from 'zustand';

type Store = {
  wsDisconnected: boolean;
  setWsDisconnected: (wsDisconnected: boolean) => void;
};

export const useStore = create<Store>((set, get) => ({
  wsDisconnected: false,
  setWsDisconnected: wsDisconnected => set({ wsDisconnected }),
}));
