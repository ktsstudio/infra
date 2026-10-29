import {
  autoUpdate,
  flip,
  offset,
  shift,
  useFloating,
} from '@floating-ui/react'

// Keep the existing menu API while using the maintained React 19 positioning engine.
export function usePopper(reference, floating, options = {}) {
  const distance =
    options.modifiers?.find(m => m.name === 'offset')?.options?.offset?.[1] ?? 0
  const { floatingStyles, update } = useFloating({
    placement: options.placement,
    elements: { reference, floating },
    middleware: [offset(distance), flip(), shift()],
    whileElementsMounted: autoUpdate,
  })
  return {
    styles: { popper: floatingStyles },
    attributes: { popper: {} },
    update,
  }
}
