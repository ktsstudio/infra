import { useEffect, useRef } from 'react'
import confetti from 'canvas-confetti'

export default function Confetti({ active, elementCount = 100 }) {
  const wasActive = useRef(false)
  useEffect(() => {
    if (active && !wasActive.current) {
      confetti({ particleCount: elementCount, disableForReducedMotion: true })
    }
    wasActive.current = active
  }, [active, elementCount])
  return null
}
