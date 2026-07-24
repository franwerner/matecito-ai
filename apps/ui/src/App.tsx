import { useEffect, useState } from 'react'

export function App() {
  const [health, setHealth] = useState('…')

  useEffect(() => {
    fetch('/api/health')
      .then((res) => res.json())
      .then((data) => setHealth(JSON.stringify(data)))
      .catch(() => setHealth('sin conexión al broker'))
  }, [])

  return (
    <main>
      <h1>Matecito UI — dev</h1>
      <p>broker /health: {health}</p>
    </main>
  )
}
