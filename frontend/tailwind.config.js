/** @type {import('tailwindcss').Config} */
export default {
    darkMode: ["class"],
    content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
  	extend: {
  		fontFamily: {
  			sans: [
  				'Inter',
  				'system-ui',
  				'sans-serif'
  			]
  		},
  		colors: {
  			/* Custom background hierarchy */
  			'bg-0': '#0a0e1a',
  			'bg-1': '#111827',
  			'bg-2': '#1e2842',
  			'bg-3': '#2a3554',
  			slate: {
  				'50': '#f8fafc',
  				'100': '#f1f5f9',
  				'200': '#e2e8f0',
  				'300': '#cbd5e1',
  				'400': '#94a3b8',
  				'500': '#64748b',
  				'600': '#475569',
  				'700': '#334155',
  				'800': '#1e293b',
  				'900': '#0f172a',
  				'950': '#020617'
  			},
  			border: 'hsl(var(--border))',
  			input: 'hsl(var(--input))',
  			ring: 'hsl(var(--ring))',
  			background: 'hsl(var(--background))',
  			foreground: 'hsl(var(--foreground))',
  			primary: {
  				DEFAULT: 'hsl(var(--primary))',
  				foreground: 'hsl(var(--primary-foreground))'
  			},
  			secondary: {
  				DEFAULT: 'hsl(var(--secondary))',
  				foreground: 'hsl(var(--secondary-foreground))'
  			},
  			destructive: {
  				DEFAULT: 'hsl(var(--destructive))',
  				foreground: 'hsl(var(--destructive-foreground))'
  			},
  			muted: {
  				DEFAULT: 'hsl(var(--muted))',
  				foreground: 'hsl(var(--muted-foreground))'
  			},
  			accent: {
  				DEFAULT: 'hsl(var(--accent))',
  				foreground: 'hsl(var(--accent-foreground))'
  			},
  			popover: {
  				DEFAULT: 'hsl(var(--popover))',
  				foreground: 'hsl(var(--popover-foreground))'
  			},
  			card: {
  				DEFAULT: 'hsl(var(--card))',
  				foreground: 'hsl(var(--card-foreground))'
  			},
  			chart: {
  				'1': 'hsl(var(--chart-1))',
  				'2': 'hsl(var(--chart-2))',
  				'3': 'hsl(var(--chart-3))',
  				'4': 'hsl(var(--chart-4))',
  				'5': 'hsl(var(--chart-5))'
  			}
  		},
  		animation: {
  			'fade-in': 'fadeIn 0.3s ease-out',
  			'slide-in': 'slideIn 0.3s ease-out',
  			'pulse-glow': 'pulse-glow 2s ease-in-out infinite',
  			'shimmer': 'shimmer 2s infinite',
  			'glow-pulse': 'glow-pulse 2s ease-in-out infinite'
  		},
  		boxShadow: {
  			'glow-primary': '0 4px 15px -3px rgba(139, 92, 246, 0.45)',
  			'glow-indigo': '0 4px 15px -3px rgba(99, 102, 241, 0.45)',
  			'glow-emerald': '0 4px 15px -3px rgba(16, 185, 129, 0.45)',
  			'glow-orange': '0 4px 15px -3px rgba(249, 115, 22, 0.45)',
  			'glow-error': '0 4px 15px -3px rgba(239, 68, 68, 0.45)',
  		},
  		keyframes: {
  			fadeIn: {
  				'0%': {
  					opacity: '0',
  					transform: 'translateY(10px)'
  				},
  				'100%': {
  					opacity: '1',
  					transform: 'translateY(0)'
  				}
  			},
  			slideIn: {
  				'0%': {
  					transform: 'translateX(-100%)'
  				},
  				'100%': {
  					transform: 'translateX(0)'
  				}
  			},
  			'pulse-glow': {
  				'0%, 100%': {
  					boxShadow: '0 0 5px rgba(99, 102, 241, 0.5)'
  				},
  				'50%': {
  					boxShadow: '0 0 20px rgba(99, 102, 241, 0.8)'
  				}
  			},
  			shimmer: {
  				'0%': { transform: 'translateX(-100%) skewX(-12deg)' },
  				'100%': { transform: 'translateX(100%) skewX(-12deg)' }
  			},
  			'glow-pulse': {
  				'0%, 100%': { opacity: '0.2' },
  				'50%': { opacity: '0.4' }
  			}
  		},
  		borderRadius: {
  			lg: 'var(--radius)',
  			md: 'calc(var(--radius) - 2px)',
  			sm: 'calc(var(--radius) - 4px)'
  		}
  	}
  },
  plugins: [require("tailwindcss-animate")],
}