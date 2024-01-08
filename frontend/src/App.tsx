import { useState } from 'react'
import { UsernameHolder } from './component/UsernameHolder'
import { ChatContainer } from './component/ChatContainer'

function App() {
  const [username, setUsername] = useState<string>("")

  const logIn = () => {
    setUsername(prompt("enter username") ?? "")
    const ws = new WebSocket("ws://localhost:8080/ws")
    ws.addEventListener("open", () => ws.send(username))
    ws.addEventListener("message", (event) => console.log(event.data))
  }

  return (
    <>
      username.length ?
      <div>
        <h4>Go Baatcheet</h4>
        <UsernameHolder sender='Placeholder - Sender' receiver='Placeholder - Sender' />
        <ChatContainer />
      </div>
      : <div>Please <button onClick={() => logIn()}>enter</button> username!</div>
    </>
  )
}

export default App
