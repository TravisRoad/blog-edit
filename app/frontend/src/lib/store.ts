import { create } from "zustand";
import { combine, devtools, persist } from "zustand/middleware";
import { useState, useEffect } from "react";

export function useStore<T, F>(
  store: (callback: (state: T) => unknown) => unknown,
  callback: (state: T) => F
) {
  const result = store(callback) as F;
  const [data, setData] = useState<F>();

  useEffect(() => {
    setData(result);
  }, [result]);

  return data;
}

const initPreference = {
  token: "",
};

const usePreference = create(
  devtools(
    persist(
      combine(initPreference, (set, get) => ({
        Set: (p: any) => {
          set(p);
        },
      })),
      { name: "user-preference" }
    )
  )
);

const useTheme = create(
  devtools(
    persist(
      combine(
        {
          theme: "",
        },
        (set, get) => ({
          Set: (s: any) => {
            set(s);
          },
        })
      ),
      { name: "theme" }
    )
  )
);

const initEditer = {
  file: "",
  content: "",
};

const useEditor = create(
  devtools(
    persist(
      combine(initEditer, (set, get) => ({
        Set: (p: any) => {
          set(p);
        },
        GetContent: (): string => {
          return get().content;
        },
      })),
      { name: "editor" }
    )
  )
);

export { usePreference, useTheme, useEditor };
