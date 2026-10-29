# Infra Dashboard

This directory contains the source code for the Infra Dashboard.

## Set up

- Install Node.js 26.13.0 or newer and npm 11.20.0 or newer.
- Run `npm ci` to install the exact, reviewed dependency versions.
- Install the following extensions if using Visual Studio Code:
  - [ESLint](https://marketplace.visualstudio.com/items?itemName=dbaeumer.vscode-eslint)
  - [Prettier](https://marketplace.visualstudio.com/items?itemName=esbenp.prettier-vscode)
  - [Tailwind CSS IntelliSense](https://marketplace.visualstudio.com/items?itemName=bradlc.vscode-tailwindcss)
  - Enable **Format on Save** in settings

## Develop

Follow the instructions in the [test config file](./__test__/__files__/infra.yaml).

```
npm ci
npm run dev
```

## Build and run

```
npm run build
npm start
```

## Linting

Linting is done via [ESLint](https://eslint.org/)

```
npm run lint
```

## Formatting

Code is formatted using [Prettier](https://prettier.io/)

To check for issues:

```
npm run format
```

To fix:

```
npm run format:fix
```

## Tests and security checks

```sh
npm test
npm audit --audit-level=low
npm outdated
```

Tests use Vitest with jsdom and Testing Library. ESLint 10 uses the maintained
Next.js and React Hooks plugins directly; the old Next.js ESLint preset includes
plugins whose peer ranges do not yet support ESLint 10. The JavaScript-only UI
has no direct TypeScript toolchain dependency.

The production build uses Next.js's supported webpack builder. This also permits
building in environments that block the internal local port used by Turbopack's
CSS worker. The standalone output and Docker entry point are unchanged.

Menu positioning uses Floating UI rather than the React 18-only react-popper
wrapper. Confetti uses canvas-confetti and respects reduced-motion preferences.
The unused react-table v7, react-router-dom and js-yaml dependencies were removed.
Jest was replaced with Vitest because the latest Jest release still depends on
a vulnerable deepmerge release without a published patched version.
