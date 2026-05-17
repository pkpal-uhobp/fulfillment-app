/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
      colors: {
        brand: {
          50: '#eff6ff',
          100: '#dbeafe',
          500: '#2563eb',
          600: '#1d4ed8',
          700: '#1e40af',
          950: '#172554',
        },
        ink: '#0f172a',
      },
      boxShadow: {
        soft: '0 24px 80px rgba(15, 23, 42, 0.10)',
      },
    },
  },
  plugins: [],
}
