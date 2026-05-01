/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        riad: {
          50: '#faf9f6',
          100: '#f0ede3',
          200: '#e1d9c9',
          300: '#d4c8a4',
          400: '#c0ad78',
          500: '#a89452',
          600: '#8a7550',
          700: '#6d5a46',
          800: '#56463a',
          900: '#3d3230',
          950: '#231c1a',
        },
        gold: {
          50: '#fefce8',
          100: '#fef9c3',
          200: '#fef08a',
          300: '#fde047',
          400: '#facc15',
          500: '#eab308',
          600: '#ca8a04',
          700: '#a16207',
          800: '#854d0e',
          900: '#713f12',
          950: '#422006',
        },
        terracotta: {
          400: '#ff9f7f',
          500: '#ff7b5a',
          600: '#e85d3a',
        }
      },
      fontFamily: {
        display: ['Georgia', 'serif'],
        arabic: ['Amiri', 'serif'],
        body: ['system-ui', '-apple-system', 'sans-serif'],
      },
      boxShadow: {
        'glass': '0 8px 32px 0 rgba(35, 28, 26, 0.1)',
        'lux': '0 10px 40px -10px rgba(168, 149, 82, 0.15)',
        'inner-gold': 'inset 0 2px 4px 0 rgba(234, 179, 8, 0.06)',
      },
      backdropBlur: {
        'xs': '2px',
      },
      animation: {
        'float': 'float 6s ease-in-out infinite',
        'pulse-slow': 'pulse 4s cubic-bezier(0.4, 0, 0.6, 1) infinite',
      },
      keyframes: {
        float: {
          '0%, 100%': { transform: 'translateY(0)' },
          '50%': { transform: 'translateY(-10px)' },
        }
      }
    },
  },
  plugins: [],
}

