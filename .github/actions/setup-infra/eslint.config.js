import js from '@eslint/js';
export default [js.configs.recommended, {ignores: ['dist/**']}, {languageOptions: {globals: {URL: 'readonly', process: 'readonly', console: 'readonly', Buffer: 'readonly', setTimeout: 'readonly'}}}];
