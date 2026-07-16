/** Tailwind preset shared between web-public and web-admin. */
/** @type {import('tailwindcss').Config} */
export default {
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        bg: {
          DEFAULT: '#060914',
          soft: '#0c1120',
          panel: '#11182b',
          'panel-2': '#161f37',
        },
        line: {
          DEFAULT: '#233252',
          soft: '#1a2540',
        },
        ink: {
          DEFAULT: '#e7edf8',
          muted: '#9dafcc',
          'muted-2': '#6c7e9c',
        },
        brand: {
          DEFAULT: '#7c3aed',
          400: '#a78bfa',
          500: '#8b5cf6',
          600: '#7c3aed',
          700: '#6d28d9',
        },
        accent: {
          DEFAULT: '#3b82f6',
        },
      },
      fontFamily: {
        sans: [
          'Inter',
          'Segoe UI',
          'Roboto',
          'system-ui',
          '-apple-system',
          'sans-serif',
        ],
      },
      boxShadow: {
        elevated: '0 18px 40px rgba(0,0,0,0.45)',
        glow: '0 0 24px rgba(124,58,237,0.35)',
      },
      borderRadius: {
        card: '14px',
      },
      keyframes: {
        'fade-in': {
          '0%': { opacity: '0', transform: 'translateY(4px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        'radar-scan': {
          '0%': { transform: 'rotate(0deg)' },
          '100%': { transform: 'rotate(360deg)' },
        },
      },
      animation: {
        'fade-in': 'fade-in .2s ease-out',
        'radar-scan': 'radar-scan 8s linear infinite',
      },
    },
  },
  plugins: [],
};
