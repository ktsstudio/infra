import js from '@eslint/js'
import { defineConfig, globalIgnores } from 'eslint/config'
import next from '@next/eslint-plugin-next'
import hooks from 'eslint-plugin-react-hooks'
import globals from 'globals'

export default defineConfig([
  js.configs.recommended,
  {
    files: ['**/*.{js,mjs,jsx}'],
    languageOptions: {
      parserOptions: { ecmaFeatures: { jsx: true } },
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.vitest,
        jsonBody: 'readonly',
      },
    },
    plugins: {
      '@next/next': next,
      'react-hooks': hooks,
      jsx: {
        rules: {
          'uses-vars': {
            create(context) {
              return {
                JSXOpeningElement(node) {
                  let name = node.name
                  while (name.type === 'JSXMemberExpression') name = name.object
                  if (
                    name.type === 'JSXIdentifier' &&
                    /^[A-Z]/.test(name.name)
                  ) {
                    context.sourceCode.markVariableAsUsed(name.name, node)
                  }
                },
              }
            },
          },
        },
      },
    },
    rules: {
      ...next.configs.recommended.rules,
      ...next.configs['core-web-vitals'].rules,
      'react-hooks/rules-of-hooks': 'error',
      '@next/next/no-img-element': 'off',
      'jsx/uses-vars': 'error',
      'no-unused-vars': [
        'error',
        { varsIgnorePattern: '^React$', args: 'after-used' },
      ],
    },
  },
  globalIgnores(['.next/**', 'out/**', 'coverage/**', 'next-env.d.ts']),
])
