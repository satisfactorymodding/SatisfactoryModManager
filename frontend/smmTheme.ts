
import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const myCustomTheme: CustomThemeConfig = {
  name: 'smm',
  properties: {
  // =~= Theme Properties =~=
    '--theme-font-family-base': '\'Roboto\', Helvetica Neue, sans-serif',
    '--theme-font-family-heading': '\'Roboto\', Helvetica Neue, sans-serif',
    '--theme-font-color-base': '0 0 0',
    '--theme-font-color-dark': '255 255 255',
    '--theme-rounded-base': '4px',
    '--theme-rounded-container': '4px',
    '--theme-border-base': '1px',
    // =~= Theme On-X Colors =~=
    '--on-primary': '0 0 0',
    '--on-secondary': '255 255 255',
    '--on-tertiary': '0 0 0',
    '--on-success': '0 0 0',
    '--on-warning': '0 0 0',
    '--on-error': 'var(--color-surface-900)',
    '--on-surface': '255 255 255',
    // =~= Theme Colors  =~=
    // primary | #66cf20 
    '--color-primary-50': '232 248 222', // #e8f8de
    '--color-primary-100': '224 245 210', // #e0f5d2
    '--color-primary-200': '217 243 199', // #d9f3c7
    '--color-primary-300': '194 236 166', // #c2eca6
    '--color-primary-400': '148 221 99', // #94dd63
    '--color-primary-500': '102 207 32', // #66cf20
    '--color-primary-600': '92 186 29', // #5cba1d
    '--color-primary-700': '77 155 24', // #4d9b18
    '--color-primary-800': '61 124 19', // #3d7c13
    '--color-primary-900': '50 101 16', // #326510
    // secondary | #666666 
    '--color-secondary-50': '232 232 232', // #e8e8e8
    '--color-secondary-100': '224 224 224', // #e0e0e0
    '--color-secondary-200': '217 217 217', // #d9d9d9
    '--color-secondary-300': '194 194 194', // #c2c2c2
    '--color-secondary-400': '148 148 148', // #949494
    '--color-secondary-500': '102 102 102', // #666666
    '--color-secondary-600': '92 92 92', // #5c5c5c
    '--color-secondary-700': '77 77 77', // #4d4d4d
    '--color-secondary-800': '61 61 61', // #3d3d3d
    '--color-secondary-900': '50 50 50', // #323232
    // tertiary | #0ea5e9 
    '--color-tertiary-50': '219 242 252', // #dbf2fc
    '--color-tertiary-100': '207 237 251', // #cfedfb
    '--color-tertiary-200': '195 233 250', // #c3e9fa
    '--color-tertiary-300': '159 219 246', // #9fdbf6
    '--color-tertiary-400': '86 192 240', // #56c0f0
    '--color-tertiary-500': '14 165 233', // #0ea5e9
    '--color-tertiary-600': '13 149 210', // #0d95d2
    '--color-tertiary-700': '11 124 175', // #0b7caf
    '--color-tertiary-800': '8 99 140', // #08638c
    '--color-tertiary-900': '7 81 114', // #075172
    // success | #66cf20 
    '--color-success-50': '232 248 222', // #e8f8de
    '--color-success-100': '224 245 210', // #e0f5d2
    '--color-success-200': '217 243 199', // #d9f3c7
    '--color-success-300': '194 236 166', // #c2eca6
    '--color-success-400': '148 221 99', // #94dd63
    '--color-success-500': '102 207 32', // #66cf20
    '--color-success-600': '92 186 29', // #5cba1d
    '--color-success-700': '77 155 24', // #4d9b18
    '--color-success-800': '61 124 19', // #3d7c13
    '--color-success-900': '50 101 16', // #326510
    // warning | #ffcc00 
    '--color-warning-50': '255 247 217', // #fff7d9
    '--color-warning-100': '255 245 204', // #fff5cc
    '--color-warning-200': '255 242 191', // #fff2bf
    '--color-warning-300': '255 235 153', // #ffeb99
    '--color-warning-400': '255 219 77', // #ffdb4d
    '--color-warning-500': '255 204 0', // #ffcc00
    '--color-warning-600': '230 184 0', // #e6b800
    '--color-warning-700': '191 153 0', // #bf9900
    '--color-warning-800': '153 122 0', // #997a00
    '--color-warning-900': '125 100 0', // #7d6400
    // error | #d60000 
    '--color-error-50': '249 217 217', // #f9d9d9
    '--color-error-100': '247 204 204', // #f7cccc
    '--color-error-200': '245 191 191', // #f5bfbf
    '--color-error-300': '239 153 153', // #ef9999
    '--color-error-400': '226 77 77', // #e24d4d
    '--color-error-500': '214 0 0', // #d60000
    '--color-error-600': '193 0 0', // #c10000
    '--color-error-700': '161 0 0', // #a10000
    '--color-error-800': '128 0 0', // #800000
    '--color-error-900': '105 0 0', // #690000
    // surface | #3b3b3b 
    '--color-surface-50': '226 226 226', // #e2e2e2
    '--color-surface-100': '216 216 216', // #d8d8d8
    '--color-surface-200': '206 206 206', // #cecece
    '--color-surface-300': '177 177 177', // #b1b1b1
    '--color-surface-400': '118 118 118', // #767676
    '--color-surface-500': '59 59 59', // #3b3b3b
    '--color-surface-600': '53 53 53', // #353535
    '--color-surface-700': '44 44 44', // #2c2c2c
    '--color-surface-800': '35 35 35', // #232323
    '--color-surface-900': '29 29 29', // #1d1d1d
  
  },
};
