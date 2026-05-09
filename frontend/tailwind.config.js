/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        ios: {
          bg: '#000000',
          card: 'rgba(255, 255, 255, 0.1)',
          text: '#FFFFFF',
          secondary: '#8E8E93',
          accent: '#0A84FF',
          danger: '#FF453A',
          warning: '#FF9F0A',
          success: '#30D158',
        }
      },
      animation: {
        'blob-pulse': 'blob-pulse 10s infinite alternate',
        'glass-glow': 'glass-glow 4s ease-in-out infinite',
      },
      keyframes: {
        'blob-pulse': {
          '0%': { transform: 'translate(0, 0) scale(1)', filter: 'blur(40px)' },
          '50%': { transform: 'translate(20px, -30px) scale(1.1)', filter: 'blur(60px)' },
          '100%': { transform: 'translate(-20px, 20px) scale(0.9)', filter: 'blur(40px)' },
        },
        'glass-glow': {
          '0%, 100%': { opacity: '0.3' },
          '50%': { opacity: '0.6' },
        }
      },
      backdropBlur: {
        'ios': '40px',
      },
      borderRadius: {
        'ios': '24px',
      }
    },
  },
  plugins: [],
}