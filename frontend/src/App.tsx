import { useState } from 'react'
import './App.css'
import { UsernameHolder } from './component/UsernameHolder'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <div>
        <h1>Chat Application</h1>
        <UsernameHolder />
        <ChatHolder />
        <ChatSender />
      </div>
    </>
  )
}

export default App
