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
      DEFAULT: 'var(--MAIN_DISABLED_BG)',
      SEARCH_RADIUS: '50px'
    },
    extend: {
      colors: {
        NEW_WARNING_COLOR: 'var(--NEW_WARNING_COLOR)',
        DISABLE_BG: 'var(--DISABLE_BG)',
        MAIN_TEXT: 'var(--MAIN_TEXT)',
        MAIN_BG: 'var(--MAIN_BG)',
        'bar-theme': 'var(--BAR_BG_COLOR)',
        BORDER: 'var(--BORDER)',
        NAVBAR_BTN_BG: 'var(--NAVBAR_BTN_BG)',
        TIPS_TEXT_COLOR: 'var(--TIPS_TEXT_COLOR)',
        LIGHT_BG_COLOR: 'var(--LIGHT_BG_COLOR)',
        SEC_TEXT: 'var(--SEC_TEXT)',
        MAIN_DISABLED_BG: 'var(--MAIN_DISABLED_BG)',
        FIX_BG: '#ffffff'
      },
      spacing: {
        mbase: 'var(--MARGIN)',
        label: '12px', // 选择器和label之间的间距，待删
        btnbase: 'var(--PADDING_BTN)', // x方向的间距
        btnybase: '16px', // y轴方向的间距
        btnrbase: '20px', // 页面最右侧边距20px
        formtop: '20px',
        icon: 'var(--MR_ICON)',
        DEFAULT_BORDER_RADIUS: 'var(--DEFAULT_BORDER_RADIUS)'
      },
      borderColor: {
        'color-base': 'var(--BORDER)'
      }
    }
  },
  plugins: []
}
