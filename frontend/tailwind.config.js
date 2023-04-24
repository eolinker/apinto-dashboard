/** @type {import('tailwindcss').Config} */
module.exports = {
  important: true,
  content: ['projects/**/*.{html,ts}'],
  theme: {
    width: {
      INPUT_NORMAL: '346px',
      INPUT_LARGE: '508px',
      GROUP: '240px',
      SEARCH: '276px'
    },
    borderRadius: {
      DEFAULT: 'var(--border-radius)',
      SEARCH_RADIUS: '50px'
    },
    extend: {
      colors: {
        NEW_WARNING_COLOR: 'var(--NEW_WARNING_COLOR)',
        DISABLE_BG: 'var(--disabled-background-color)',
        MAIN_TEXT: 'var(--text-color)',
        MAIN_BG: 'var(--background-color)',
        'bar-theme': 'var(--bar-background-color)',
        BORDER: 'var(--border-color)',
        NAVBAR_BTN_BG: 'var(--item-active-background-color)',
        TIPS_TEXT_COLOR: 'var(--text-secondary-color)',
        LIGHT_BG_COLOR: 'var(--LIGHT_BG_COLOR)',
        SEC_TEXT: 'var(--SEC_TEXT)',
        MAIN_DISABLED_BG: 'var(--disabled-background-color)',
        FIX_BG: '#ffffff',
        theme: 'var(--primary-color)',
        DESC_TEXT: '#666666',
        HOVER_BG: 'var(--item-hover-background-color)'
      },
      spacing: {
        mbase: '20px',
        label: '12px', // 选择器和label之间的间距，待删
        btnbase: '10px', // x方向的间距
        btnybase: '10px', // y轴方向的间距
        btnrbase: '20px', // 页面最右侧边距20px
        formtop: '20px',
        icon: '5px',
        DEFAULT_BORDER_RADIUS: 'var(--border-radius)'
      },
      borderColor: {
        'color-base': 'var(--border-color)'
      }
    }
  },
  plugins: []
}
