import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  plugins: [],
  theme: {
    extend: {
      fontFamily: {
        mono: "Robot Mono",
      },
    },
  },
  safelist: ["text-amber-400", "text-green-400", "text-red-400"],
};
export default config;
